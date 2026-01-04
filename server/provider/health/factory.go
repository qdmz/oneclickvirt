package health

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// CreateHealthChecker 根据provider配置创建健康检查器的便捷函数
func CreateHealthChecker(providerType, host, username, password string, port int, logger *zap.Logger) (HealthChecker, error) {
	config := HealthConfig{
		Host:       host,
		Port:       port,
		Username:   username,
		Password:   password,
		SSHEnabled: true,
		APIEnabled: true,
	}

	manager := NewHealthManager(logger)
	return manager.CreateChecker(ProviderType(providerType), config)
}

// DefaultHealthConfig 创建默认的健康检查配置
func DefaultHealthConfig() HealthConfig {
	return HealthConfig{
		Port:           22,
		SSHEnabled:     true,
		APIEnabled:     true,
		APIScheme:      "https",
		Timeout:        30000000000, // 30 seconds in nanoseconds
		ServiceChecks:  []string{},
		CustomCommands: []string{},
	}
}

// HealthCheckAdapter 适配器，用于将新的健康检查系统与现有代码集成
type HealthCheckAdapter struct {
	checker HealthChecker
}

// NewHealthCheckAdapter 创建健康检查适配器
func NewHealthCheckAdapter(checker HealthChecker) *HealthCheckAdapter {
	return &HealthCheckAdapter{
		checker: checker,
	}
}

// CheckHealth 执行健康检查并返回简化结果
func (hca *HealthCheckAdapter) CheckHealth(ctx context.Context) error {
	result, err := hca.checker.CheckHealth(ctx)
	if err != nil {
		return err
	}

	if result.Status == HealthStatusUnhealthy {
		if len(result.Errors) > 0 {
			return fmt.Errorf("health check failed: %s", result.Errors[0])
		}
		return fmt.Errorf("health check failed")
	}

	return nil
}

// GetHealthResult 获取详细的健康检查结果
func (hca *HealthCheckAdapter) GetHealthResult(ctx context.Context) (*HealthResult, error) {
	return hca.checker.CheckHealth(ctx)
}
