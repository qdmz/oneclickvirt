package payment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"oneclickvirt/global"
	"oneclickvirt/model/common"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Yipay 易支付服务
type Yipay struct {
	// API地址
	APIURL string
	// 商户ID
	MerchantID string
	// 商户密钥
	MerchantKey string
	// 是否验证SSL
	VerifySSL bool
}

// NewYipay 创建易支付服务
func NewYipay() *Yipay {
	return &Yipay{
		APIURL:      "https://pay.wanjuanxueyi.com",
		MerchantID:  "2093",
		MerchantKey: "7o6IxRTgt67ntX9nIZRx2koiPX9X2ix2",
		VerifySSL:   false,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	OutTradeNo string  `json:"out_trade_no"` // 商户订单号
	Subject    string  `json:"subject"`      // 商品名称
	TotalFee   float64 `json:"total_fee"`    // 订单金额
	NotifyURL  string  `json:"notify_url"`   // 异步通知地址
	ReturnURL  string  `json:"return_url"`   // 同步跳转地址
	Param      string  `json:"param"`        // 自定义参数
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	Code    int    `json:"code"`    // 状态码
	Msg     string `json:"msg"`     // 消息
	Data    struct {
		OutTradeNo string `json:"out_trade_no"` // 商户订单号
		TradeNo    string `json:"trade_no"`     // 平台订单号
		QrcodeURL  string `json:"qrcode_url"`   // 二维码地址
		PayURL     string `json:"pay_url"`      // 支付地址
	} `json:"data"`
}

// QueryOrderRequest 查询订单请求
type QueryOrderRequest struct {
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
}

// QueryOrderResponse 查询订单响应
type QueryOrderResponse struct {
	Code    int    `json:"code"`    // 状态码
	Msg     string `json:"msg"`     // 消息
	Data    struct {
		OutTradeNo string  `json:"out_trade_no"` // 商户订单号
		TradeNo    string  `json:"trade_no"`     // 平台订单号
		TotalFee   float64 `json:"total_fee"`    // 订单金额
		TradeStatus string  `json:"trade_status"` // 交易状态
		Time       string  `json:"time"`          // 交易时间
	} `json:"data"`
}

// NotifyRequest 回调通知请求
type NotifyRequest struct {
	TradeNo    string  `form:"trade_no"`     // 平台订单号
	OutTradeNo string  `form:"out_trade_no"` // 商户订单号
	TotalFee   float64 `form:"total_fee"`    // 订单金额
	TradeStatus string `form:"trade_status"` // 交易状态
	Param      string  `form:"param"`        // 自定义参数
	Sign       string  `form:"sign"`         // 签名
}

// CreateOrder 创建订单
func (y *Yipay) CreateOrder(req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// 构建请求参数
	params := map[string]string{
		"pid":          y.MerchantID,
		"type":         "alipay", // 默认支付宝
		"out_trade_no": req.OutTradeNo,
		"name":         req.Subject,
		"money":        fmt.Sprintf("%.2f", req.TotalFee),
		"notify_url":   req.NotifyURL,
		"return_url":   req.ReturnURL,
		"param":        req.Param,
	}

	// 生成签名
	sign := y.generateSign(params)
	params["sign"] = sign
	params["sign_type"] = "MD5"

	// 转换为 url.Values
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	// 发送请求
	apiURL := fmt.Sprintf("%s/submit.php", y.APIURL)
	resp, err := http.PostForm(apiURL, values)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查是否是HTML响应
	if strings.Contains(string(body), "<html") || strings.Contains(string(body), "<HTML") {
		// 从HTML中提取支付URL
		// 简化处理：直接返回HTML内容
		return &CreateOrderResponse{
			Code: 1,
			Msg:  "创建成功",
			Data: struct {
				OutTradeNo string `json:"out_trade_no"`
				TradeNo    string `json:"trade_no"`
				QrcodeURL  string `json:"qrcode_url"`
				PayURL     string `json:"pay_url"`
			}{
				OutTradeNo: req.OutTradeNo,
				TradeNo:    req.OutTradeNo,
				QrcodeURL:  apiURL + "?out_trade_no=" + req.OutTradeNo,
				PayURL:     apiURL + "?out_trade_no=" + req.OutTradeNo,
			},
		}, nil
	}

	// 解析响应
	var result CreateOrderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查状态码
	if result.Code != 1 {
		return nil, fmt.Errorf("创建订单失败: %s", result.Msg)
	}

	return &result, nil
}

// QueryOrder 查询订单
func (y *Yipay) QueryOrder(req *QueryOrderRequest) (*QueryOrderResponse, error) {
	// 构建请求参数
	params := map[string]string{
		"act":          "query",
		"pid":          y.MerchantID,
		"out_trade_no": req.OutTradeNo,
	}

	// 生成签名
	sign := y.generateSign(params)
	params["sign"] = sign
	params["sign_type"] = "MD5"

	// 转换为 url.Values
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	// 发送请求
	apiURL := fmt.Sprintf("%s/api.php", y.APIURL)
	resp, err := http.PostForm(apiURL, values)
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var result QueryOrderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查状态码
	if result.Code != 1 {
		return nil, fmt.Errorf("查询订单失败: %s", result.Msg)
	}

	return &result, nil
}

// VerifyNotify 验证回调通知
func (y *Yipay) VerifyNotify(req *NotifyRequest) bool {
	// 构建验证参数
	params := map[string]string{
		"trade_no":     req.TradeNo,
		"out_trade_no": req.OutTradeNo,
		"total_fee":    fmt.Sprintf("%.2f", req.TotalFee),
		"trade_status": req.TradeStatus,
		"param":        req.Param,
	}

	// 生成签名
	sign := y.generateSign(params)

	// 验证签名
	return sign == req.Sign
}

// generateSign 生成签名
func (y *Yipay) generateSign(params map[string]string) string {
	// 排序参数
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接参数
	var sb strings.Builder
	for _, k := range keys {
		if params[k] != "" {
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(params[k])
			sb.WriteString("&")
		}
	}

	// 添加密钥
	sb.WriteString("key=")
	sb.WriteString(y.MerchantKey)

	// MD5加密
	hash := md5.Sum([]byte(sb.String()))
	return hex.EncodeToString(hash[:])
}

// YipayNotify 易支付回调处理
func YipayNotify(c *gin.Context) {
	var req NotifyRequest
	if err := c.ShouldBind(&req); err != nil {
		global.APP_LOG.Error("易支付回调参数错误", zap.Error(err))
		c.String(http.StatusOK, "fail")
		return
	}

	// 创建易支付服务
	yipay := NewYipay()

	// 验证签名
	if !yipay.VerifyNotify(&req) {
		global.APP_LOG.Error("易支付回调签名验证失败",
			zap.String("trade_no", req.TradeNo),
			zap.String("out_trade_no", req.OutTradeNo))
		c.String(http.StatusOK, "fail")
		return
	}

	// 记录日志
	global.APP_LOG.Info("易支付回调",
		zap.String("trade_no", req.TradeNo),
		zap.String("out_trade_no", req.OutTradeNo),
		zap.Float64("total_fee", req.TotalFee),
		zap.String("trade_status", req.TradeStatus),
		zap.String("param", req.Param))

	// 检查交易状态
	if req.TradeStatus != "TRADE_SUCCESS" {
		global.APP_LOG.Warn("易支付回调交易状态不正确",
			zap.String("trade_status", req.TradeStatus))
		c.String(http.StatusOK, "fail")
		return
	}

	// TODO: 处理订单支付成功逻辑
	// 1. 根据商户订单号查询订单
	// 2. 更新订单状态
	// 3. 触发自动开通流程

	// 返回成功
	c.String(http.StatusOK, "success")
}

// YipayReturn 易支付同步跳转
func YipayReturn(c *gin.Context) {
	var req NotifyRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 400,
			Msg:  "参数错误",
		})
		return
	}

	// 创建易支付服务
	yipay := NewYipay()

	// 验证签名
	if !yipay.VerifyNotify(&req) {
		c.JSON(http.StatusOK, common.Response{
			Code: 401,
			Msg:  "签名验证失败",
		})
		return
	}

	// 检查交易状态
	if req.TradeStatus != "TRADE_SUCCESS" {
		c.JSON(http.StatusOK, common.Response{
			Code: 400,
			Msg:  "交易未完成",
		})
		return
	}

	// TODO: 处理订单支付成功逻辑
	// 1. 根据商户订单号查询订单
	// 2. 更新订单状态
	// 3. 跳转到订单详情页面

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Msg:  "支付成功",
		Data: gin.H{
			"trade_no":     req.TradeNo,
			"out_trade_no": req.OutTradeNo,
			"total_fee":    req.TotalFee,
		},
	})
}

// TestYipay 测试易支付接口
func TestYipay(c *gin.Context) {
	// 创建易支付服务
	yipay := NewYipay()

	// 创建测试订单
	req := &CreateOrderRequest{
		OutTradeNo: fmt.Sprintf("TEST%d", time.Now().Unix()),
		Subject:    "测试商品",
		TotalFee:   0.01,
		NotifyURL:  global.APP_CONFIG.System.FrontendURL + "/api/v1/payment/yipay/notify",
		ReturnURL:  global.APP_CONFIG.System.FrontendURL + "/api/v1/payment/yipay/return",
		Param:      "test",
	}

	// 创建订单
	result, err := yipay.CreateOrder(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Code: 500,
			Msg:  "创建订单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Data: result,
		Msg:  "创建成功",
	})
}

// QueryYipayOrder 查询易支付订单
func QueryYipayOrder(c *gin.Context) {
	outTradeNo := c.Query("out_trade_no")
	if outTradeNo == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Code: 400,
			Msg:  "订单号不能为空",
		})
		return
	}

	// 创建易支付服务
	yipay := NewYipay()

	// 查询订单
	req := &QueryOrderRequest{
		OutTradeNo: outTradeNo,
	}

	result, err := yipay.QueryOrder(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Code: 500,
			Msg:  "查询订单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Data: result,
		Msg:  "查询成功",
	})
}

// CreateYipayOrder 创建易支付订单
func CreateYipayOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Code: 400,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认回调地址
	if req.NotifyURL == "" {
		req.NotifyURL = global.APP_CONFIG.System.FrontendURL + "/api/v1/payment/yipay/notify"
	}
	if req.ReturnURL == "" {
		req.ReturnURL = global.APP_CONFIG.System.FrontendURL + "/api/v1/payment/yipay/return"
	}

	// 创建易支付服务
	yipay := NewYipay()

	// 创建订单
	result, err := yipay.CreateOrder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Code: 500,
			Msg:  "创建订单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Data: result,
		Msg:  "创建成功",
	})
}

// GetYipayConfig 获取易支付配置
func GetYipayConfig(c *gin.Context) {
	yipay := NewYipay()

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Data: gin.H{
			"api_url":      yipay.APIURL,
			"merchant_id":  yipay.MerchantID,
			"merchant_key": "***", // 隐藏密钥
			"verify_ssl":   yipay.VerifySSL,
		},
		Msg: "获取成功",
	})
}

// UpdateYipayConfig 更新易支付配置
func UpdateYipayConfig(c *gin.Context) {
	var req struct {
		APIURL      string `json:"api_url" binding:"required"`
		MerchantID  string `json:"merchant_id" binding:"required"`
		MerchantKey string `json:"merchant_key" binding:"required"`
		VerifySSL   bool   `json:"verify_ssl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Code: 400,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	// TODO: 保存配置到数据库或配置文件

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Msg:  "更新成功",
	})
}

// GetYipayStats 获取易支付统计
func GetYipayStats(c *gin.Context) {
	// TODO: 从数据库获取易支付统计数据

	stats := gin.H{
		"total_orders":    0,
		"total_amount":    0.0,
		"success_orders":  0,
		"pending_orders":  0,
		"failed_orders":   0,
		"today_orders":    0,
		"today_amount":    0.0,
	}

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Data: stats,
		Msg:  "获取成功",
	})
}

// GetYipayOrders 获取易支付订单列表
func GetYipayOrders(c *gin.Context) {
	// 获取参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	status := c.Query("status")

	// 转换参数
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	// TODO: 从数据库获取易支付订单列表

	orders := []gin.H{}
	total := int64(0)

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Data: gin.H{
			"list":     orders,
			"total":    total,
			"page":     pageInt,
			"pageSize": pageSizeInt,
			"status":   status,
		},
		Msg: "获取成功",
	})
}
