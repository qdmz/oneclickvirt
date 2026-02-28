package payment

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"oneclickvirt/global"
	orderModel "oneclickvirt/model/order"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 易支付相关配置
type EpayConfig struct {
	APIURL    string `json:"api_url"`
	PID       string `json:"pid"`
	Key       string `json:"key"`
	ReturnURL string `json:"return_url"`
	NotifyURL string `json:"notify_url"`
	Type      string `json:"type"` // alipay, wechat, qqpay
}

// 码支付相关配置
type MapayConfig struct {
	APIURL    string `json:"api_url"`
	ID        string `json:"id"`
	Key       string `json:"key"`
	ReturnURL string `json:"return_url"`
	NotifyURL string `json:"notify_url"`
	Type      string `json:"type"` // alipay, wechat, qqpay
}

// EpayNotify 易支付回调
// @Summary 易支付回调
// @Description 易支付异步通知
// @Tags 支付
// @Accept json
// @Produce json
// @Router /v1/payment/epay/notify [post]
func EpayNotify(c *gin.Context) {
	// 解析表单数据
	if err := c.Request.ParseForm(); err != nil {
		global.APP_LOG.Error("解析易支付回调数据失败", zap.Error(err))
		c.String(http.StatusBadRequest, "fail")
		return
	}

	params := c.Request.Form
	sign := params.Get("sign")
	
	// 验证签名
	if !verifyEpaySign(params, sign) {
		global.APP_LOG.Error("易支付签名验证失败")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 获取订单号
	orderNo := params.Get("out_trade_no")
	if orderNo == "" {
		global.APP_LOG.Error("易支付订单号为空")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 处理支付成功
	if err := processPaymentSuccess(orderNo, orderModel.PaymentMethodEpay, valuesToMap(params)); err != nil {
		global.APP_LOG.Error("处理易支付成功失败", zap.Error(err))
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// MapayNotify 码支付回调
// @Summary 码支付回调
// @Description 码支付异步通知
// @Tags 支付
// @Accept json
// @Produce json
// @Router /v1/payment/mapay/notify [post]
func MapayNotify(c *gin.Context) {
	// 解析表单数据
	if err := c.Request.ParseForm(); err != nil {
		global.APP_LOG.Error("解析码支付回调数据失败", zap.Error(err))
		c.String(http.StatusBadRequest, "fail")
		return
	}

	params := c.Request.Form
	sign := params.Get("sign")
	
	// 验证签名
	if !verifyMapaySign(params, sign) {
		global.APP_LOG.Error("码支付签名验证失败")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 获取订单号
	orderNo := params.Get("out_trade_no")
	if orderNo == "" {
		global.APP_LOG.Error("码支付订单号为空")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 处理支付成功
	if err := processPaymentSuccess(orderNo, orderModel.PaymentMethodMapay, valuesToMap(params)); err != nil {
		global.APP_LOG.Error("处理码支付成功失败", zap.Error(err))
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// verifyEpaySign 验证易支付签名
func verifyEpaySign(params url.Values, sign string) bool {
	// 构建待签名字符串
	var keys []string
	for key := range params {
		if key != "sign" && key != "sign_type" && params.Get(key) != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var signStr strings.Builder
	for i, key := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(key)
		signStr.WriteString("=")
		signStr.WriteString(params.Get(key))
	}
	signStr.WriteString(global.APP_CONFIG.Payment.EpayKey)

	// MD5签名
	h := md5.New()
	h.Write([]byte(signStr.String()))
	calculatedSign := hex.EncodeToString(h.Sum(nil))

	return strings.ToLower(calculatedSign) == strings.ToLower(sign)
}

// verifyMapaySign 验证码支付签名
func verifyMapaySign(params url.Values, sign string) bool {
	// 构建待签名字符串
	var keys []string
	for key := range params {
		if key != "sign" && key != "sign_type" && params.Get(key) != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var signStr strings.Builder
	for i, key := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(key)
		signStr.WriteString("=")
		signStr.WriteString(params.Get(key))
	}
	signStr.WriteString(global.APP_CONFIG.Payment.MapayKey)

	// MD5签名
	h := md5.New()
	h.Write([]byte(signStr.String()))
	calculatedSign := hex.EncodeToString(h.Sum(nil))

	return strings.ToLower(calculatedSign) == strings.ToLower(sign)
}

// 将url.Values转换为map[string]interface{}
func valuesToMap(values url.Values) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range values {
		if len(value) > 0 {
			result[key] = value[0]
		}
	}
	return result
}