package site

import (
	"time"
)

// SiteConfig 站点配置表
type SiteConfig struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Key         string    `json:"key" gorm:"type:varchar(100);uniqueIndex;not null;comment:配置键"`
	Value       string    `json:"value" gorm:"type:text;comment:配置值"`
	Type        string    `json:"type" gorm:"type:varchar(20);not null;comment:配置类型"`
	Description string    `json:"description" gorm:"type:varchar(255);comment:配置说明"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (SiteConfig) TableName() string {
	return "site_configs"
}

// 预定义的站点配置键常量
const (
	SiteIconURL   = "site_icon_url"  // 网站图标URL
	SiteName      = "site_name"      // 网站名称
	SiteURL       = "site_url"       // 网站URL
	SiteHeader    = "site_header"    // 页眉内容
	SiteFooter    = "site_footer"    // 页脚内容
	CustomCSS     = "custom_css"     // 自定义CSS
	ContactEmail  = "contact_email"  // 联系邮箱
	ContactPhone  = "contact_phone"  // 联系电话
	CompanyName   = "company_name"   // 公司名称
	ICPNumber     = "icp_number"     // ICP备案号
	AnalyticsCode = "analytics_code" // 统计代码
)

// 配置类型常量
const (
	ConfigTypeString = "string" // 字符串类型
	ConfigTypeImage  = "image"  // 图片类型
	ConfigTypeJSON   = "json"   // JSON类型
	ConfigTypeText   = "text"   // 文本类型
)
