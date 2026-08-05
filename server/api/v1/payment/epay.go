package payment

import (
	"oneclickvirt/global"
	orderSvc "oneclickvirt/service/product"
	walletSvc "oneclickvirt/service/wallet"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EpayNotify 易支付异步回调（公开，无需鉴权）
// 同一笔交易可能既是充值也可能是购买，根据 out_trade_no 前缀区分：
// RC 开头 = 余额充值；PD 开头 = 产品购买。
// @Router /payment/epay/notify [post]
func EpayNotify(c *gin.Context) {
	params := map[string]string{}
	if err := c.Request.ParseForm(); err == nil {
		for k, v := range c.Request.PostForm {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	}
	// 易支付也支持 GET 回调（部分渠道），兼容 query 参数
	if len(params) == 0 {
		for k, v := range c.Request.URL.Query() {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	}

	outTradeNo := params["out_trade_no"]
	if outTradeNo == "" {
		c.String(200, "fail")
		return
	}

	var handleErr error
	switch {
	case len(outTradeNo) >= 2 && outTradeNo[:2] == "RC":
		handleErr = walletSvc.NewService().HandleEpayNotify(params)
	case len(outTradeNo) >= 2 && outTradeNo[:2] == "PD":
		handleErr = orderSvc.NewService().HandleEpayNotify(params)
	default:
		c.String(200, "fail")
		return
	}

	if handleErr != nil {
		global.APP_LOG.Error("易支付回调处理失败", zap.Error(handleErr))
		c.String(200, "fail")
		return
	}
	// 易支付要求返回 success 表示已接收
	c.String(200, "success")
}

// EpayReturn 易支付同步跳转（公开）
// @Router /payment/epay/return [get]
func EpayReturn(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "支付完成，正在返回...",
		"data": gin.H{
			"out_trade_no": c.Query("out_trade_no"),
		},
	})
}
