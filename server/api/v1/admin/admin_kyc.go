package admin

import (
	"strconv"

	"oneclickvirt/global"
	kycModel "oneclickvirt/model/kyc"
	"oneclickvirt/model/common"
	"oneclickvirt/service/kyc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetKYCRecords handles GET /v1/admin/kyc/records
func GetKYCRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filters := make(map[string]interface{})
	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filters["status"] = status
		}
	}
	if username := c.Query("username"); username != "" {
		filters["username"] = username
	}

	svc, err := getAdminKYCService()
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeExternalAPIError, "KYC service unavailable"))
		return
	}

	records, total, err := svc.GetAllKYCRecords(page, pageSize, filters)
	if err != nil {
		global.APP_LOG.Error("get KYC records failed", zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "failed to query records"))
		return
	}

	// Format records
	list := make([]gin.H, 0, len(records))
	for i := range records {
		r := &records[i]
		item := gin.H{
			"id":          r.ID,
			"userId":      r.UserID,
			"certifyId":   r.CertifyID,
			"realName":    kycModel.MaskRealName(r.RealName),
			"idType":      r.IDType,
			"gender":      r.Gender,
			"province":    r.Province,
			"city":        r.City,
			"status":      r.Status,
			"statusText":  statusText(r.Status),
			"certifiedAt": r.CertifiedAt,
			"createdAt":   r.CreatedAt,
			"remark":      r.Remark,
		}
		// Include username if preloaded
		if r.UserID != 0 {
			item["userId"] = r.UserID
		}
		list = append(list, item)
	}

	common.ResponseSuccessWithPagination(c, list, total, page, pageSize)
}

// UpdateKYCStatus handles PUT /v1/admin/kyc/:id/status
func UpdateKYCStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "invalid record ID"))
		return
	}

	var req struct {
		Status int    `json:"status" binding:"required"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, err.Error()))
		return
	}

	svc, err := getAdminKYCService()
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeExternalAPIError, "KYC service unavailable"))
		return
	}

	if err := svc.UpdateKYCStatus(uint(id), req.Status, req.Remark); err != nil {
		global.APP_LOG.Error("update KYC status failed", zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "status updated")
}

// GetKYCStats handles GET /v1/admin/kyc/stats
func GetKYCStats(c *gin.Context) {
	db := global.APP_DB
	var total, pending, certified, rejected, expired int64

	db.Model(&kycModel.KYCRecord{}).Count(&total)
	db.Model(&kycModel.KYCRecord{}).Where("status = ?", kycModel.KYCStatusPending).Count(&pending)
	db.Model(&kycModel.KYCRecord{}).Where("status = ?", kycModel.KYCStatusCertified).Count(&certified)
	db.Model(&kycModel.KYCRecord{}).Where("status = ?", kycModel.KYCStatusRejected).Count(&rejected)
	db.Model(&kycModel.KYCRecord{}).Where("status = ?", kycModel.KYCStatusExpired).Count(&expired)

	common.ResponseSuccess(c, gin.H{
		"total":     total,
		"pending":   pending,
		"certified": certified,
		"rejected":  rejected,
		"expired":   expired,
	})
}

func getAdminKYCService() (*kyc.KYCService, error) {
	paymentCfg := global.APP_CONFIG.Payment
	return kyc.NewKYCService(global.APP_DB, &kyc.PaymentConfig{
		EnableRealName:      paymentCfg.EnableRealName,
		RequireRealName:     paymentCfg.RequireRealName,
		AlipayAppID:         paymentCfg.AlipayAppID,
		AlipayPrivateKey:    paymentCfg.AlipayPrivateKey,
		AlipayPublicKey:     paymentCfg.AlipayPublicKey,
		AlipayGateway:       paymentCfg.AlipayGateway,
		RealNameCallbackURL: paymentCfg.RealNameCallbackURL,
	})
}

func statusText(status int) string {
	switch status {
	case kycModel.KYCStatusPending:
		return "pending"
	case kycModel.KYCStatusCertified:
		return "certified"
	case kycModel.KYCStatusRejected:
		return "rejected"
	case kycModel.KYCStatusExpired:
		return "expired"
	default:
		return "unknown"
	}
}
