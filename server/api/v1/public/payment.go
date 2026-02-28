package public

import (
	"net/http"
	"oneclickvirt/global"
	"oneclickvirt/model/common"

	"github.com/gin-gonic/gin"
)

// GetPaymentConfig 获取支付配置（公开接口）
// @Summary 获取支付配置
// @Description 获取启用的支付方式配置（不包含敏感信息）
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Router /v1/public/payment-config [get]
func GetPaymentConfig(c *gin.Context) {
	config := global.APP_CONFIG.Payment

	// 只返回支付方式是否启用，不包含敏感配置信息
	paymentConfig := gin.H{
		"enableAlipay":  config.EnableAlipay,
		"enableWechat":  config.EnableWechat,
		"enableBalance": config.EnableBalance,
		"enableEpay":    config.EnableEpay,
		"enableMapay":   config.EnableMapay,
	}

	c.JSON(http.StatusOK, common.Success(paymentConfig))
}
