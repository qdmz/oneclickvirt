package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"oneclickvirt/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Viper 初始化配置文件
func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		config = "config.yaml"
	} else {
		config = path[0]
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")

	// 支持从环境变量中读取配置
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("[VIPER] 配置文件读取错误: %s，使用默认配置\n", err)
		// 不要panic，而是使用默认配置继续运行
		return v
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("[VIPER] 配置文件已更改: %s\n", e.Name)
		if err := v.Unmarshal(&global.APP_CONFIG); err != nil {
			fmt.Printf("[VIPER] 配置解析失败: %v\n", err)
		}
	})

	if err := v.Unmarshal(&global.APP_CONFIG); err != nil {
		fmt.Printf("[VIPER] 配置解析失败: %v\n", err)
	}

	// 设置默认值
	setDefaults(v)

	return v
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("system.env", "public")
	v.SetDefault("system.addr", 8080)
	v.SetDefault("system.db-type", "mysql")
	v.SetDefault("system.oss-type", "local")
	v.SetDefault("system.use-multipoint", false)
	v.SetDefault("system.use-redis", false)
	v.SetDefault("system.iplimit-count", 15000)
	v.SetDefault("system.iplimit-time", 3600)

	// 数据库默认配置
	v.SetDefault("mysql.path", "127.0.0.1")
	v.SetDefault("mysql.port", "3306")
	v.SetDefault("mysql.config", "charset=utf8mb4&parseTime=True&loc=Local")
	v.SetDefault("mysql.db-name", "oneclickvirt")
	v.SetDefault("mysql.username", "root")
	v.SetDefault("mysql.password", "root")
	v.SetDefault("mysql.prefix", "")
	v.SetDefault("mysql.singular", false)
	v.SetDefault("mysql.engine", "InnoDB")
	v.SetDefault("mysql.max-idle-conns", 10)
	v.SetDefault("mysql.max-open-conns", 100)
	v.SetDefault("mysql.log-mode", "info")
	v.SetDefault("mysql.log-zap", false)
	v.SetDefault("mysql.max-lifetime", 3600)
	v.SetDefault("mysql.auto-create", true)

	// 生成强制的安全JWT签名密钥
	randomKey := generateSecureJWTKey()

	// 验证密钥强度
	if err := validateJWTKeyStrength(randomKey); err != nil {
		fmt.Printf("[VIPER] JWT密钥强度验证失败: %v，重新生成密钥\n", err)
		// 重新生成密钥
		randomKey = generateSecureJWTKey()
	}

	v.SetDefault("jwt.signing-key", randomKey)
	v.SetDefault("jwt.expires-time", "7d")
	v.SetDefault("jwt.buffer-time", "1d")
	v.SetDefault("jwt.issuer", "oneclickvirt")

	v.SetDefault("zap.level", "info")
	v.SetDefault("zap.format", "console")
	v.SetDefault("zap.prefix", "[oneclickvirt]")
	v.SetDefault("zap.director", "logs")
	v.SetDefault("zap.show-line", true)
	v.SetDefault("zap.encode-level", "LowercaseColorLevelEncoder")
	v.SetDefault("zap.stacktrace-key", "stacktrace")
	v.SetDefault("zap.log-in-console", true)
}

// generateSecureJWTKey 生成安全的JWT密钥
func generateSecureJWTKey() string {
	bytes := make([]byte, 32) // 强制256位密钥
	if _, err := rand.Read(bytes); err != nil {
		// 如果生成失败，使用时间戳作为后备，但确保足够长度
		backupKey := fmt.Sprintf("oneclickvirt-backup-%d", time.Now().UnixNano())
		// 确保至少64字符长度
		for len(backupKey) < 64 {
			backupKey += fmt.Sprintf("-%d", time.Now().UnixNano())
		}
		return backupKey[:64] // 截取到64字符
	}
	return hex.EncodeToString(bytes)
}

// validateJWTKeyStrength 验证JWT密钥强度
func validateJWTKeyStrength(key string) error {
	if len(key) < 32 {
		return fmt.Errorf("JWT密钥长度不足，当前长度: %d，最小要求: 32", len(key))
	}

	// 检查是否是弱密钥
	weakKeys := []string{
		"secret",
		"password",
		"12345",
		"test",
		"jwt-secret",
		"your-secret-key",
		"change-me",
	}

	for _, weak := range weakKeys {
		if strings.Contains(strings.ToLower(key), weak) {
			return fmt.Errorf("JWT密钥包含弱模式，请使用更强的密钥")
		}
	}

	return nil
}
