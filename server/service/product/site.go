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
	"oneclickvirt/utils"

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

	// 获取易支付启用的支付方式
	yipayEnabledPayTypes := ""
	var yipayConfig product.YiPayConfig
	if err := global.APP_DB.Where("enabled = ?", true).First(&yipayConfig).Error; err == nil {
		yipayEnabledPayTypes = yipayConfig.EnabledPayTypes
		if yipayEnabledPayTypes == "" {
			yipayEnabledPayTypes = "alipay,wxpay,qqpay"
		}
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
		"show_platforms":       config.ShowPlatformsSection,
		"show_sponsors":        config.ShowSponsorsSection,
		"show_recommended":     config.ShowRecommendedSection,
		"recommended_limit":    config.RecommendedLimit,
		"recommended_title":    config.RecommendedTitle,
		"recommended_subtitle": config.RecommendedSubtitle,
		"recommended_cols":     config.RecommendedCols,
		"recommended_show_price": config.RecommendedShowPrice,
		"recommended_show_specs": config.RecommendedShowSpecs,
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
		"yipay_pay_types":      yipayEnabledPayTypes,
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

// allowedSiteConfigKeys 站点配置允许通过接口更新的字段白名单（JSON key == 数据库列名）
var allowedSiteConfigKeys = map[string]bool{
	"site_name": true, "site_description": true, "site_keywords": true,
	"logo_url": true, "favicon_url": true, "dark_logo_url": true,
	"custom_header": true, "header_enabled": true, "show_nav": true,
	"custom_footer": true, "footer_enabled": true, "copyright_text": true, "icp_number": true,
	"home_title": true, "home_subtitle": true, "home_background": true, "show_home_stats": true,
	"primary_color": true, "theme_mode": true, "custom_css": true,
	"contact_email": true, "contact_phone": true, "contact_qq": true, "contact_telegram": true,
	"show_balance": true, "show_yipay": true,
	"enable_registration": true, "enable_ticket": true, "enable_product_store": true,
	"announcement_bar": true, "announcement_enabled": true,
	"show_platforms": true, "show_sponsors": true, "show_recommended": true, "recommended_limit": true,
	"recommended_title": true, "recommended_subtitle": true, "recommended_cols": true,
	"recommended_show_price": true, "recommended_show_specs": true,
}

// jsonKeyToColumn 修正「JSON 字段名与数据库列名不一致」的字段。
// 模型字段 FooterHTML / HeaderHTML 的 JSON 名为 custom_footer / custom_header，
// 但 gorm 默认列名为 footer_html / header_html；若直接以 JSON 名作为更新列名会报 Unknown column。
var jsonKeyToColumn = map[string]string{
	"custom_header": "header_html",
	"custom_footer": "footer_html",
	"show_yipay":    "show_yi_pay",
}

// UpdateSiteConfigFields 部分更新站点配置
// 只更新请求中提供的字段（白名单内），未提供的字段保持不变，
// 从而修复此前使用 Save() 整体覆盖导致未提交字段（如首页配置）被清零的问题。
func (s *SiteService) UpdateSiteConfigFields(m map[string]interface{}) error {
	if len(m) == 0 {
		return nil
	}

	updates := make(map[string]interface{})
	for k, v := range m {
		if !allowedSiteConfigKeys[k] {
			continue
		}
		// 数值型字段（如 recommended_limit / recommended_cols）在 JSON 反序列化后为 float64，转回 int 避免类型不匹配
		if k == "recommended_limit" || k == "recommended_cols" {
			if f, ok := v.(float64); ok {
				v = int(f)
			}
		}
		// custom_header / custom_footer 会在公开页以 v-html 原样渲染，保存前
		// 必须过服务端 HTML 净化（剥 script/iframe/on* 事件与危险 URL），
		// 防止管理员录入的 HTML 变成面向全体访客的存储型 XSS。
		v = sanitizeSiteConfigValue(k, v)
		// JSON 字段名与数据库列名可能不一致（如 custom_footer -> footer_html），更新时需转换为真实列名
		col := k
		if c, ok := jsonKeyToColumn[k]; ok {
			col = c
		}
		updates[col] = v
	}
	if len(updates) == 0 {
		return nil
	}

	// 确保配置行存在（GetSiteConfig 会在缺失时创建默认行）
	if _, err := s.GetSiteConfig(); err != nil {
		return err
	}

	var existing product.SiteConfig
	if err := global.APP_DB.First(&existing).Error; err != nil {
		return err
	}

	if err := global.APP_DB.Model(&existing).Updates(updates).Error; err != nil {
		global.APP_LOG.Error("更新站点配置失败", zap.Error(err))
		return err
	}

	s.invalidateCache()
	global.APP_LOG.Info("更新站点配置成功", zap.Uint("configID", existing.ID), zap.Int("fields", len(updates)))
	return nil
}

// sanitizeSiteConfigValue 对会被公开页原样渲染的站点 HTML 字段做保存前净化；
// 其他字段原样返回。
func sanitizeSiteConfigValue(key string, v interface{}) interface{} {
	if key != "custom_header" && key != "custom_footer" {
		return v
	}
	if s, ok := v.(string); ok {
		return utils.SanitizeHTML(s)
	}
	return v
}

// UploadSiteImage 上传站点图片并返回URL
func (s *SiteService) UploadSiteImage(imageData []byte, filename string) (string, error) {
	if len(imageData) == 0 {
		return "", errors.New("图片数据不能为空")
	}

	// 验证文件类型。SVG 不在白名单内：SVG 可内嵌脚本，同源提供时有脚本
	// 执行面（P2-5），如确需 SVG 应先做 XML 净化再放行，当前直接移除。
	validExts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true}
	ext := strings.ToLower(filepath.Ext(filename))
	if !validExts[ext] {
		return "", errors.New("不支持的图片格式，仅支持 png, jpg, jpeg, gif, webp")
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
