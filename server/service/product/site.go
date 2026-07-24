package product

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"oneclickvirt/global"
	"oneclickvirt/model/product"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SiteService 站点配置服务
type SiteService struct {
	cache     *product.SiteConfig
	cacheTime time.Time
	mu        sync.RWMutex
	cacheTTL  time.Duration
}

var (
	siteServiceInstance *SiteService
	siteServiceOnce     sync.Once
)

// NewSiteService 创建站点配置服务实例（单例）
func NewSiteService() *SiteService {
	siteServiceOnce.Do(func() {
		siteServiceInstance = &SiteService{
			cacheTTL: 5 * time.Minute,
		}
	})
	return siteServiceInstance
}

// getCachedConfig 获取缓存的配置
func (s *SiteService) getCachedConfig() *product.SiteConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.cache != nil && time.Since(s.cacheTime) < s.cacheTTL {
		return s.cache
	}
	return nil
}

// setCachedConfig 设置缓存
func (s *SiteService) setCachedConfig(config *product.SiteConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = config
	s.cacheTime = time.Now()
}

// invalidateCache 使缓存失效
func (s *SiteService) invalidateCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = nil
	s.cacheTime = time.Time{}
}

// GetSiteConfig 获取站点配置（公开，带缓存）
func (s *SiteService) GetSiteConfig() (*product.SiteConfig, error) {
	// 先尝试从缓存获取
	if cached := s.getCachedConfig(); cached != nil {
		return cached, nil
	}

	var config product.SiteConfig
	if err := global.APP_DB.First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有配置时返回默认配置
			config = product.SiteConfig{
			SiteName:           "OneClickVirt",
			PrimaryColor:       "#409EFF",
			ThemeMode:          "auto",
			ShowNav:            true,
			ShowHomeStats:      true,
			EnableRegistration: true,
			EnableTicket:       true,
			EnableProductStore: true,
			ShowBalance:        true,
		}
			// 创建默认配置
			if err := global.APP_DB.Create(&config).Error; err != nil {
				global.APP_LOG.Error("创建默认站点配置失败", zap.Error(err))
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	s.setCachedConfig(&config)
	return &config, nil
}

// GetPublicSiteInfo 获取公开站点信息（过滤敏感字段）
func (s *SiteService) GetPublicSiteInfo() (map[string]interface{}, error) {
	config, err := s.GetSiteConfig()
	if err != nil {
		return nil, err
	}

	// 只返回前端需要的公开字段
	return map[string]interface{}{
		"site_name":            config.SiteName,
		"site_description":     config.SiteDescription,
		"site_keywords":        config.SiteKeywords,
		"logo_url":             config.LogoURL,
		"favicon_url":          config.FaviconURL,
		"dark_logo_url":        config.DarkLogoURL,
		"home_title":           config.HomeTitle,
		"home_subtitle":        config.HomeSubtitle,
		"home_background":      config.HomeBackground,
		"show_home_stats":      config.ShowHomeStats,
		"primary_color":        config.PrimaryColor,
		"theme_mode":           config.ThemeMode,
		"custom_css":           config.CustomCSS,
		"show_nav":             config.ShowNav,
		"copyright_text":       config.CopyrightText,
		"icp_number":           config.ICPNumber,
		"contact_email":        config.ContactEmail,
		"contact_phone":        config.ContactPhone,
		"contact_qq":           config.ContactQQ,
		"contact_telegram":     config.ContactTelegram,
		"show_balance":         config.ShowBalance,
		"show_yipay":           config.ShowYiPay,
		"enable_registration":  config.EnableRegistration,
		"enable_ticket":        config.EnableTicket,
		"enable_product_store": config.EnableProductStore,
		"announcement_bar":     config.AnnouncementBar,
		"announcement_enabled": config.AnnouncementEnabled,
		"custom_header":        config.HeaderHTML,
		"custom_footer":        config.FooterHTML,
		"header_enabled":       config.HeaderEnabled,
		"footer_enabled":       config.FooterEnabled,
	}, nil
}

// GetAdminSiteConfig 获取完整站点配置（管理员）
func (s *SiteService) GetAdminSiteConfig() (*product.SiteConfig, error) {
	return s.GetSiteConfig()
}

// UpdateSiteConfig 更新站点配置
func (s *SiteService) UpdateSiteConfig(config *product.SiteConfig) error {
	if config == nil {
		return errors.New("配置不能为空")
	}

	var existing product.SiteConfig
	if err := global.APP_DB.First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			if err := global.APP_DB.Create(config).Error; err != nil {
				global.APP_LOG.Error("创建站点配置失败", zap.Error(err))
				return err
			}
			s.setCachedConfig(config)
			global.APP_LOG.Info("创建站点配置成功")
			return nil
		}
		return err
	}

	// 更新现有配置
	config.ID = existing.ID
	if err := global.APP_DB.Model(&existing).Updates(config).Error; err != nil {
		global.APP_LOG.Error("更新站点配置失败", zap.Error(err))
		return err
	}

	// 使缓存失效
	s.invalidateCache()

	global.APP_LOG.Info("更新站点配置成功", zap.Uint("configID", existing.ID))
	return nil
}

// UploadSiteImage 上传站点图片并返回URL
func (s *SiteService) UploadSiteImage(imageData []byte, filename string) (string, error) {
	if len(imageData) == 0 {
		return "", errors.New("图片数据不能为空")
	}

	// 验证文件类型
	validExts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".svg": true, ".webp": true}
	ext := strings.ToLower(filepath.Ext(filename))
	if !validExts[ext] {
		return "", errors.New("不支持的图片格式，仅支持 png, jpg, jpeg, gif, svg, webp")
	}

	// 限制文件大小（5MB）
	const maxSize = 5 * 1024 * 1024
	if len(imageData) > maxSize {
		return "", errors.New("图片大小超过5MB限制")
	}

	// 构建存储路径
	uploadDir := filepath.Join("uploads", "site")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		global.APP_LOG.Error("创建上传目录失败", zap.String("dir", uploadDir), zap.Error(err))
		return "", err
	}

	// 生成唯一文件名
	uniqueName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, uniqueName)

	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		global.APP_LOG.Error("保存图片失败", zap.String("path", filePath), zap.Error(err))
		return "", err
	}

	// 返回可访问的URL路径
	imageURL := fmt.Sprintf("/uploads/site/%s", uniqueName)
	global.APP_LOG.Info("上传站点图片成功", zap.String("url", imageURL), zap.String("originalName", filename))
	return imageURL, nil
}
