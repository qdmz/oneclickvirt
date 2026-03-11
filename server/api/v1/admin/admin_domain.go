package admin

import (
	"strconv"

	"oneclickvirt/global"
	"oneclickvirt/model/common"
	"oneclickvirt/service/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getDomainService() *domain.DomainService {
	return domain.NewDomainService(global.APP_DB)
}

// AdminGetDomains 管理员获取所有域名
func AdminGetDomains(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	filters := make(map[string]interface{})
	if v := c.Query("userId"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filters["userId"] = uint(id)
		}
	}
	if v := c.Query("agentId"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			filters["agentId"] = uint(id)
		}
	}
	if v := c.Query("domain"); v != "" {
		filters["domain"] = v
	}
	if v := c.Query("status"); v != "" {
		if status, err := strconv.Atoi(v); err == nil {
			filters["status"] = status
		}
	}

	svc := getDomainService()
	domains, total, err := svc.AdminGetDomains(page, pageSize, filters)
	if err != nil {
		global.APP_LOG.Error("管理员获取域名列表失败", zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "获取域名列表失败"))
		return
	}

	common.ResponseSuccessWithPagination(c, domains, total, page, pageSize)
}

// AdminDeleteDomain 管理员删除域名
func AdminDeleteDomain(c *gin.Context) {
	domainID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "域名ID格式错误"))
		return
	}

	svc := getDomainService()
	if err := svc.AdminDeleteDomain(uint(domainID)); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil)
}

// GetDomainConfig 获取域名配置
func GetDomainConfig(c *gin.Context) {
	svc := getDomainService()
	config, err := svc.GetDomainConfig()
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "获取域名配置失败"))
		return
	}

	common.ResponseSuccess(c, config)
}

// UpdateDomainConfig 更新域名配置
func UpdateDomainConfig(c *gin.Context) {
	var req struct {
		MaxDomainsPerUser      int    `json:"maxDomainsPerUser"`
		MaxDomainsPerAgentUser int    `json:"maxDomainsPerAgentUser"`
		DefaultTTL             int    `json:"defaultTTL"`
		AutoSSL                bool   `json:"autoSSL"`
		AllowedSuffixes        string `json:"allowedSuffixes"`
		DNSType                string `json:"dnsType"`
		DNSConfigPath          string `json:"dnsConfigPath"`
		NginxConfigPath        string `json:"nginxConfigPath"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "参数错误"))
		return
	}

	svc := getDomainService()
	config, err := svc.GetDomainConfig()
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "获取当前配置失败"))
		return
	}

	config.MaxDomainsPerUser = req.MaxDomainsPerUser
	config.MaxDomainsPerAgentUser = req.MaxDomainsPerAgentUser
	config.DefaultTTL = req.DefaultTTL
	config.AutoSSL = req.AutoSSL
	config.AllowedSuffixes = req.AllowedSuffixes
	config.DNSType = req.DNSType
	config.DNSConfigPath = req.DNSConfigPath
	config.NginxConfigPath = req.NginxConfigPath

	if err := svc.UpdateDomainConfig(config); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "更新配置失败"))
		return
	}

	common.ResponseSuccess(c, config)
}

// SyncDNS 手动同步DNS
func SyncDNS(c *gin.Context) {
	svc := getDomainService()
	if err := svc.SyncAllDNS(); err != nil {
		global.APP_LOG.Error("DNS同步失败", zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "DNS同步失败: "+err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "DNS同步完成")
}
