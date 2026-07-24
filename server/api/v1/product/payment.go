package product

import (
	"net/url"
	"strconv"

	"oneclickvirt/global"
	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	productModel "oneclickvirt/model/product"
	productService "oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateYiPayOrder 创建易支付订单
// @Summary 创建易支付充值订单
// @Description 创建易支付充值订单，返回支付跳转链接
// @Tags 充值中心
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body product.CreateYiPayOrderRequest true "创建易支付订单请求参数"
// @Success 200 {object} common.Response{data=object} "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Router /payments/yipay [post]
func CreateYiPayOrder(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req productModel.CreateYiPayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	clientIP := c.ClientIP()
	service := productService.NewService()
	result, err := service.CreateYiPayOrder(userID, req, clientIP)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccess(c, result, "支付订单创建成功")
}

// YiPayNotify 易支付异步通知回调
// @Summary 易支付异步通知回调
// @Description 易支付异步通知回调接口，用于处理支付结果
// @Tags 充值中心
// @Accept x-www-form-urlencoded
// @Produce json
// @Param pid formData string true "商户ID"
// @Param trade_no formData string true "平台订单号"
// @Param out_trade_no formData string true "商户订单号"
// @Param type formData string true "支付方式"
// @Param name formData string true "商品名称"
// @Param money formData string true "订单金额"
// @Param trade_status formData string true "交易状态"
// @Param sign formData string true "签名"
// @Param sign_type formData string true "签名类型"
// @Success 200 {string} string "success"
// @Failure 400 {string} string "fail"
// @Router /payments/yipay/notify [post]
func YiPayNotify(c *gin.Context) {
	// 解析所有表单参数
	params := make(map[string]string)
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	// 如果 PostForm 为空，尝试解析 RawQuery（某些网关可能使用 query 参数）
	if len(params) == 0 {
		query, _ := url.ParseQuery(c.Request.URL.RawQuery)
		for key, values := range query {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
	}

	global.APP_LOG.Info("收到易支付异步通知",
		zap.Any("params", params))

	service := productService.NewService()
	if err := service.ProcessYiPayNotify(params); err != nil {
		global.APP_LOG.Error("易支付通知处理失败", zap.Error(err))
		c.String(200, "fail")
		return
	}

	c.String(200, "success")
}

// YiPayReturn 易支付同步跳转
// @Summary 易支付同步跳转
// @Description 易支付支付完成后同步跳转回前端页面
// @Tags 充值中心
// @Accept json
// @Produce json
// @Param out_trade_no query string true "商户订单号"
// @Success 200 {object} common.Response "处理成功"
// @Router /payments/yipay/return [get]
func YiPayReturn(c *gin.Context) {
	outTradeNo := c.Query("out_trade_no")
	if outTradeNo == "" {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "缺少订单号参数"))
		return
	}

	// 同步跳转一般只做查询展示，实际状态以异步通知为准
	common.ResponseSuccess(c, map[string]interface{}{
		"outTradeNo": outTradeNo,
		"message":    "支付结果处理中，请稍后查询余额",
	}, "支付跳转成功")
}

// GetRechargeList 获取充值记录
// @Summary 获取充值记录
// @Description 获取当前登录用户的易支付充值记录
// @Tags 充值中心
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=common.PageResult} "获取成功"
// @Failure 401 {object} common.Response "用户未登录"
// @Router /payments/recharge-records [get]
func GetRechargeList(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	service := productService.NewService()
	logs, total, err := service.GetRechargeList(userID, page, pageSize)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}

	common.ResponseSuccessWithPagination(c, logs, total, page, pageSize)
}
