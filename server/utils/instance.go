package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// GenerateInstanceName 生成实例名称（全局统一函数）
// 生成格式: provider-name-4位随机字符 (如: docker-d73a)
func GenerateInstanceName(providerName string) string {
	n, err := rand.Int(rand.Reader, big.NewInt(65536))
	if err != nil {
		n = big.NewInt(0) // fallback
	}
	randomStr := fmt.Sprintf("%04x", n.Int64()) // 生成4位16进制随机字符

	// 清理provider名称，移除特殊字符
	cleanName := strings.ReplaceAll(strings.ToLower(providerName), " ", "-")
	cleanName = strings.ReplaceAll(cleanName, "_", "-")

	return fmt.Sprintf("%s-%s", cleanName, randomStr)
}
