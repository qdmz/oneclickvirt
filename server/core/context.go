package core

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"oneclickvirt/config"
)

// AppContext 应用上下文，用于减少全局变量的使用
type AppContext struct {
	mu        sync.RWMutex
	db        *gorm.DB
	logger    *zap.Logger
	config    *config.Server
	viper     *viper.Viper
	ginEngine *gin.Engine
	configMgr *config.ConfigManager
}

var (
	appCtx *AppContext
	once   sync.Once
)

// GetAppContext 获取应用上下文实例
func GetAppContext() *AppContext {
	once.Do(func() {
		appCtx = &AppContext{}
	})
	return appCtx
}

// InitAppContext 初始化应用上下文
func InitAppContext(db *gorm.DB, logger *zap.Logger, cfg *config.Server, v *viper.Viper) {
	ctx := GetAppContext()
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.db = db
	ctx.logger = logger
	ctx.config = cfg
	ctx.viper = v

	// 初始化配置管理器
	ctx.configMgr = config.NewConfigManager(db, logger)
}

// DB 获取数据库连接
func (ctx *AppContext) DB() *gorm.DB {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.db
}

// Logger 获取日志器
func (ctx *AppContext) Logger() *zap.Logger {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.logger
}

// Config 获取配置
func (ctx *AppContext) Config() *config.Server {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.config
}

// Viper 获取Viper实例
func (ctx *AppContext) Viper() *viper.Viper {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.viper
}

// GinEngine 获取Gin引擎
func (ctx *AppContext) GinEngine() *gin.Engine {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.ginEngine
}

// SetGinEngine 设置Gin引擎
func (ctx *AppContext) SetGinEngine(engine *gin.Engine) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.ginEngine = engine
}

// ConfigManager 获取配置管理器
func (ctx *AppContext) ConfigManager() *config.ConfigManager {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.configMgr
}

// WithContext 创建带有应用上下文的context.Context
func (ctx *AppContext) WithContext(parent context.Context) context.Context {
	return context.WithValue(parent, "app_context", ctx)
}

// FromContext 从context.Context中获取应用上下文
func FromContext(ctx context.Context) *AppContext {
	if appCtx, ok := ctx.Value("app_context").(*AppContext); ok {
		return appCtx
	}
	return GetAppContext()
}

// Shutdown 优雅关闭
func (ctx *AppContext) Shutdown() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	// 关闭数据库连接
	if ctx.db != nil {
		if sqlDB, err := ctx.db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	// 同步日志
	if ctx.logger != nil {
		ctx.logger.Sync()
	}

	return nil
}

// IsInitialized 检查是否已初始化
func (ctx *AppContext) IsInitialized() bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.db != nil && ctx.logger != nil && ctx.config != nil
}
