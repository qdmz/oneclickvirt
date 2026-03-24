package kyc

import (
	"crypto/sha256"
	"fmt"
	"time"

	"oneclickvirt/global"
	kycModel "oneclickvirt/model/kyc"
	userModel "oneclickvirt/model/user"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// KYCService handles KYC (real name verification) operations
type KYCService struct {
	db        *gorm.DB
	alipayKYC *AlipayKYC
	config    *PaymentConfig
}

// PaymentConfig holds alipay config for KYC
type PaymentConfig struct {
	EnableRealName    bool
	RequireRealName   bool
	AlipayAppID       string
	AlipayPrivateKey  string
	AlipayPublicKey   string
	AlipayGateway     string
	RealNameCallbackURL string
}

// NewKYCService creates a new KYCService
func NewKYCService(db *gorm.DB, cfg *PaymentConfig) (*KYCService, error) {
	svc := &KYCService{
		db:     db,
		config: cfg,
	}
	if cfg.AlipayAppID != "" && cfg.AlipayPrivateKey != "" {
		kyc, err := NewAlipayKYC(
			cfg.AlipayAppID,
			cfg.AlipayPrivateKey,
			cfg.AlipayPublicKey,
			cfg.AlipayGateway,
			cfg.RealNameCallbackURL,
		)
		if err != nil {
			global.APP_LOG.Error("failed to init alipay KYC client", zap.Error(err))
			return svc, err
		}
		svc.alipayKYC = kyc
	}
	return svc, nil
}

// SubmitCertification submits name+ID for verification
func (s *KYCService) SubmitCertification(userID uint, realName, idCardNumber string) (*kycModel.KYCRecord, string, error) {
	// Check if already certified
	var existing kycModel.KYCRecord
	if err := s.db.Where("user_id = ? AND status = ?", userID, kycModel.KYCStatusCertified).First(&existing).Error; err == nil {
		return &existing, "", fmt.Errorf("already certified")
	}

	// Compute ID card hash for dedup
	hash := sha256.Sum256([]byte(idCardNumber))
	hashStr := fmt.Sprintf("%x", hash)

	// Check dedup
	var dup kycModel.KYCRecord
	if err := s.db.Where("id_card_hash = ? AND status = ? AND deleted_at IS NULL", hashStr, kycModel.KYCStatusCertified).First(&dup).Error; err == nil {
		return nil, "", fmt.Errorf("this ID card has already been verified")
	}

	// Check if alipay KYC is initialized
	if s.alipayKYC == nil {
		return nil, "", fmt.Errorf("KYC service not initialized")
	}

	// Call alipay
	certifyID, certifyURL, err := s.alipayKYC.InitializeCertify(realName, idCardNumber, "IDENTITY_CARD")
	if err != nil {
		return nil, "", fmt.Errorf("alipay initialize failed: %v", err)
	}

	// Delete old pending records for this user
	s.db.Where("user_id = ? AND status = ?", userID, kycModel.KYCStatusPending).Delete(&kycModel.KYCRecord{})

	// Create record
	record := &kycModel.KYCRecord{
		UserID:        userID,
		CertifyID:     certifyID,
		RealName:      realName,
		IDCardNumber:  idCardNumber,
		IDCardHash:    hashStr,
		Status:        kycModel.KYCStatusPending,
		IDType:        "IDENTITY_CARD",
	}
	if err := s.db.Create(record).Error; err != nil {
		return nil, "", fmt.Errorf("create KYC record: %v", err)
	}

	return record, certifyURL, nil
}

// CallbackCertify handles alipay callback
func (s *KYCService) CallbackCertify(certifyID string) error {
	var record kycModel.KYCRecord
	if err := s.db.Where("certify_id = ?", certifyID).First(&record).Error; err != nil {
		return fmt.Errorf("KYC record not found for certify_id: %s", certifyID)
	}

	if record.Status != kycModel.KYCStatusPending {
		return nil // already processed
	}

	passed, _, err := s.alipayKYC.QueryCertify(certifyID)
	if err != nil {
		return fmt.Errorf("query certify: %v", err)
	}

	now := time.Now()
	if passed {
		record.Status = kycModel.KYCStatusCertified
		record.CertifiedAt = &now
		// Update user verification flag
		s.db.Model(&userModel.User{}).Where("id = ?", record.UserID).Update("real_name_verified", true)
	} else {
		record.Status = kycModel.KYCStatusRejected
		record.Remark = "verification failed"
	}
	return s.db.Save(&record).Error
}

// GetKYCStatus gets the user's KYC status
func (s *KYCService) GetKYCStatus(userID uint) (*kycModel.KYCRecord, error) {
	var record kycModel.KYCRecord
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").First(&record).Error
	if err != nil {
		return nil, nil // no record
	}
	return &record, nil
}

// QueryAndUpdateCertification manually queries and updates certification status
func (s *KYCService) QueryAndUpdateCertification(userID uint) (*kycModel.KYCRecord, error) {
	var record kycModel.KYCRecord
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").First(&record).Error; err != nil {
		return nil, fmt.Errorf("no KYC record found")
	}

	if record.Status != kycModel.KYCStatusPending {
		return &record, nil
	}

	passed, _, err := s.alipayKYC.QueryCertify(record.CertifyID)
	if err != nil {
		return nil, fmt.Errorf("query certify: %v", err)
	}

	now := time.Now()
	if passed {
		record.Status = kycModel.KYCStatusCertified
		record.CertifiedAt = &now
		s.db.Model(&userModel.User{}).Where("id = ?", record.UserID).Update("real_name_verified", true)
	} else {
		record.Status = kycModel.KYCStatusRejected
		record.Remark = "verification failed"
	}
	s.db.Save(&record)
	return &record, nil
}

// GetAllKYCRecords returns paginated KYC records for admin
func (s *KYCService) GetAllKYCRecords(page, pageSize int, filters map[string]interface{}) ([]kycModel.KYCRecord, int64, error) {
	var records []kycModel.KYCRecord
	var total int64

	query := s.db.Model(&kycModel.KYCRecord{})

	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if username, ok := filters["username"]; ok {
		query = query.Joins("JOIN users ON users.id = kyc_records.user_id").
			Where("users.username LIKE ?", "%"+fmt.Sprintf("%v", username)+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).
		Find(&records).Error
	return records, total, err
}

// UpdateKYCStatus allows admin to override KYC status
func (s *KYCService) UpdateKYCStatus(recordID uint, status int, remark string) error {
	var record kycModel.KYCRecord
	if err := s.db.First(&record, recordID).Error; err != nil {
		return fmt.Errorf("record not found")
	}

	record.Status = status
	record.Remark = remark
	if status == kycModel.KYCStatusCertified && record.CertifiedAt == nil {
		now := time.Now()
		record.CertifiedAt = &now
		s.db.Model(&userModel.User{}).Where("id = ?", record.UserID).Update("real_name_verified", true)
	}
	return s.db.Save(&record).Error
}

// GetPendingKYCCount returns the count of pending KYC records
func (s *KYCService) GetPendingKYCCount() (int64, error) {
	var count int64
	err := s.db.Model(&kycModel.KYCRecord{}).Where("status = ?", kycModel.KYCStatusPending).Count(&count).Error
	return count, err
}
