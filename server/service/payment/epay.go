package payment

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"oneclickvirt/global"
)

// EpaySign 计算易支付签名。
// 规则：将除 sign/sign_type 外的参数按 key 升序排序，拼接为 k1=v1&k2=v2...，
// 末尾追加通信密钥 key，整体做 md5（小写十六进制）。
func EpaySign(params map[string]string, key string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if v == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
	}
	sb.WriteString(key)

	sum := md5.Sum([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// EpayVerify 校验易支付异步回调签名。
func EpayVerify(params map[string]string, key string) bool {
	sign, ok := params["sign"]
	if !ok || sign == "" {
		return false
	}
	return strings.EqualFold(sign, EpaySign(params, key))
}

// BuildEpaySubmitParams 构造跳转易支付收银台所需参数（含签名）。
// money 为元（形如 "10.00"），name 为商品名，outTradeNo 为系统订单号。
func BuildEpaySubmitParams(outTradeNo, name, money, notifyURL, returnURL string) (map[string]string, error) {
	cfg := global.GetAppConfig().Payment
	if cfg.EpayPID == "" || cfg.EpayKey == "" || cfg.EpayApiURL == "" {
		return nil, fmt.Errorf("易支付未配置（请在后台填写 PID/KEY/网关地址）")
	}
	params := map[string]string{
		"pid":        cfg.EpayPID,
		"type":       "alipay", // 默认支付宝，前端/后台可扩展
		"out_trade_no": outTradeNo,
		"notify_url": notifyURL,
		"return_url": returnURL,
		"name":       name,
		"money":      money,
		"sign_type":  "MD5",
	}
	params["sign"] = EpaySign(params, cfg.EpayKey)
	return params, nil
}

// EpaySubmitURL 返回易支付下单网关地址。
func EpaySubmitURL() string {
	return strings.TrimRight(global.GetAppConfig().Payment.EpayApiURL, "/") + "/submit.php"
}

// GenerateTradeNo 生成系统内部订单号（out_trade_no）。
func GenerateTradeNo(prefix string) string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return prefix + timeNowUnix() + hex.EncodeToString(b)
}
