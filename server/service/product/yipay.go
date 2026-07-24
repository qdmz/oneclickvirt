package product

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"oneclickvirt/global"
	"oneclickvirt/model/common"
	productModel "oneclickvirt/model/product"

	"go.uber.org/zap"
)

// YiPayService 易支付服务
type YiPayService struct{}

// NewYiPayService 创建易支付服务实例
func NewYiPayService() *YiPayService {
	return &YiPayService{}
}

// GetActiveConfig 获取启用的易支付配置
func (s *YiPayService) GetActiveConfig() (*productModel.YiPayConfig, error) {
	var config productModel.YiPayConfig
	if err := global.APP_DB.Where("enabled = ?", true).First(&config).Error; err != nil {
		return nil, common.NewError(common.CodeNotFound, "易支付未配置或未启用")
	}
	return &config, nil
}

// GenerateOrderNo 生成充值订单号
// 格式: RCH + 年月日 + 6位随机数
func (s *YiPayService) GenerateOrderNo() string {
	now := time.Now()
	randomNum := fmt.Sprintf("%06d", randInt(0, 999999))
	return "RCH" + now.Format("20060102") + randomNum
}

// BuildPayParams 构建易支付请求参数
func (s *YiPayService) BuildPayParams(config *productModel.YiPayConfig, orderNo string, amount float64, payType string, clientIP string) map[string]string {
	params := map[string]string{
		"pid":          config.Pid,
		"type":         payType,
		"out_trade_no": orderNo,
		"notify_url":   config.NotifyURL,
		"return_url":   config.ReturnURL,
		"name":         "账户充值",
		"money":        fmt.Sprintf("%.2f", amount),
		"clientip":     clientIP,
		"device":       "pc",
	}

	// 生成签名
	params["sign"] = s.Sign(params, config.Key)
	params["sign_type"] = "MD5"

	return params
}

// GeneratePayURL 生成支付跳转URL
func (s *YiPayService) GeneratePayURL(config *productModel.YiPayConfig, params map[string]string) string {
	payURL, _ := url.Parse(config.ApiURL + "/submit.php")
	query := payURL.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	payURL.RawQuery = query.Encode()
	return payURL.String()
}

// VerifyNotify 验证易支付异步通知
func (s *YiPayService) VerifyNotify(params map[string]string, config *productModel.YiPayConfig) bool {
	sign := params["sign"]
	if sign == "" {
		global.APP_LOG.Warn("易支付通知缺少签名")
		return false
	}

	// 移除签名相关字段后验证
	delete(params, "sign")
	delete(params, "sign_type")

	expectedSign := s.Sign(params, config.Key)
	if !strings.EqualFold(sign, expectedSign) {
		global.APP_LOG.Warn("易支付通知签名验证失败",
			zap.String("received", sign),
			zap.String("expected", expectedSign))
		return false
	}

	return true
}

// Sign 生成易支付MD5签名
// 规则: 参数按ASCII码排序后拼接成URL键值对格式，末尾加上key，再进行MD5加密
func (s *YiPayService) Sign(params map[string]string, key string) string {
	// 过滤空值和签名字段
	filtered := make(map[string]string)
	for k, v := range params {
		if v == "" || k == "sign" || k == "sign_type" {
			continue
		}
		filtered[k] = v
	}

	// 按key排序
	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接字符串
	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(filtered[k])
	}
	sb.WriteString(key)

	// MD5加密并转大写
	hash := md5.Sum([]byte(sb.String()))
	return fmt.Sprintf("%032X", hash)
}

// QueryOrder 查询易支付订单状态
func (s *YiPayService) QueryOrder(config *productModel.YiPayConfig, orderNo string) (map[string]interface{}, error) {
	params := map[string]string{
		"pid":          config.Pid,
		"out_trade_no": orderNo,
	}
	params["sign"] = s.Sign(params, config.Key)
	params["sign_type"] = "MD5"

	queryURL, _ := url.Parse(config.ApiURL + "/api.php")
	q := queryURL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	queryURL.RawQuery = q.Encode()

	resp, err := http.Get(queryURL.String())
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 简单解析JSON（实际项目中建议使用json.Unmarshal到结构体）
	global.APP_LOG.Debug("易支付订单查询响应", zap.String("body", string(body)))

	// 这里返回原始数据，由调用方解析
	return map[string]interface{}{
		"raw": string(body),
	}, nil
}

// randInt 生成指定范围的随机整数
func randInt(min, max int) int {
	if min >= max {
		return min
	}
	// 使用当前时间纳秒作为简单随机源
	n := int(time.Now().UnixNano())
	return min + (n % (max - min + 1))
}

// Float64ToString float64转字符串，保留2位小数
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
