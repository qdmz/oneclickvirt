package payment

import (
	"strconv"
	"time"
)

// timeNowUnix 返回当前 Unix 时间戳字符串（避免直接依赖 time 在 epay.go 中引入多余 import）。
func timeNowUnix() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// FormatMoney 将分转换为易支付所需的元字符串（保留两位）。
func FormatMoney(cents int64) string {
	return strconv.FormatFloat(float64(cents)/100.0, 'f', 2, 64)
}

// ParseMoneyCents 将元字符串解析为分（int64）。
func ParseMoneyCents(money string) (int64, error) {
	f, err := strconv.ParseFloat(money, 64)
	if err != nil {
		return 0, err
	}
	return int64(f * 100 + 0.5), nil
}
