package user

import (
	"regexp"

	"oneclickvirt/global"
	kycModel "oneclickvirt/model/kyc"
	"oneclickvirt/model/common"
	"oneclickvirt/service/kyc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// isValidIDCard validates Chinese ID card number (18 digits)
func isValidIDCard(id string) bool {
	if len(id) != 18 {
		return false
	}
	matched, _ := regexp.MatchString(`^\d{17}[\dXx]$`, id)
	return matched
}

// SubmitCertification handles POST /v1/user/kyc/submit
func SubmitCertification(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req struct {
		RealName     string `json:"realName" binding:"required"`
		IDCardNumber string `json:"idCardNumber" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, err.Error()))
		return
	}

	if !isValidIDCard(req.IDCardNumber) {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "invalid ID card number"))
		return
	}

	paymentCfg := global.APP_CONFIG.Payment
	if !paymentCfg.EnableRealName {
		common.ResponseWithError(c, common.NewError(common.CodeError, "real name verification is not enabled"))
		return
	}

	svc, err := getKYCService()
	if err != nil {
		global.APP_LOG.Error("init KYC service failed", zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeExternalAPIError, "KYC service unavailable"))
		return
	}

	record, certifyURL, err := svc.SubmitCertification(userID, req.RealName, req.IDCardNumber)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeError, err.Error()))
		return
	}

	common.ResponseSuccess(c, gin.H{
		"record":     formatKYCRecord(record),
		"certifyURL": certifyURL,
	})
}

// GetKYCStatus handles GET /v1/user/kyc/status
func GetKYCStatus(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	svc, err := getKYCService()
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeExternalAPIError, "KYC service unavailable"))
		return
	}

	record, err := svc.GetKYCStatus(userID)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	if record == nil {
		common.ResponseSuccess(c, gin.H{"status": -1, "message": "not verified"})
		return
	}

	common.ResponseSuccess(c, formatKYCRecord(record))
}

// QueryAndUpdate handles POST /v1/user/kyc/query
func QueryAndUpdate(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	svc, err := getKYCService()
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeExternalAPIError, "KYC service unavailable"))
		return
	}

	record, err := svc.QueryAndUpdateCertification(userID)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeError, err.Error()))
		return
	}

	common.ResponseSuccess(c, formatKYCRecord(record))
}

// HandleKYCCallback handles GET /v1/kyc/callback (public route)
func HandleKYCCallback(c *gin.Context) {
	certifyID := c.Query("certify_id")
	if certifyID == "" {
		c.String(400, "missing certify_id")
		return
	}

	svc, err := getKYCService()
	if err != nil {
		global.APP_LOG.Error("KYC callback service init failed", zap.Error(err))
		c.String(500, "service unavailable")
		return
	}

	if err := svc.CallbackCertify(certifyID); err != nil {
		global.APP_LOG.Error("KYC callback failed", zap.Error(err))
		c.String(500, "callback failed")
		return
	}

	frontendURL := global.APP_CONFIG.System.FrontendURL
	if frontendURL == "" {
		frontendURL = "/"
	}
	c.Redirect(302, frontendURL+"/#/user/kyc?callback=success")
}

func getKYCService() (*kyc.KYCService, error) {
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

func formatKYCRecord(r *kycModel.KYCRecord) gin.H {
	return gin.H{
		"id":          r.ID,
		"userId":      r.UserID,
		"certifyId":   r.CertifyID,
		"realName":    kycModel.MaskRealName(r.RealName),
		"status":      r.Status,
		"statusText":  statusText(r.Status),
		"certifiedAt": r.CertifiedAt,
		"idType":      r.IDType,
		"gender":      r.Gender,
		"createdAt":   r.CreatedAt,
		"remark":      r.Remark,
	}
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
