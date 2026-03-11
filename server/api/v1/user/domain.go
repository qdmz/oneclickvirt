package user

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

// GetDomains 获取用户域名列表
func GetDomains(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	svc := getDomainService()
	domains, total, err := svc.GetDomains(userID, page, pageSize)
	if err != nil {
		global.APP_LOG.Error("获取域名列表失败", zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "获取域名列表失败"))
		return
	}

	common.ResponseSuccessWithPagination(c, domains, total, page, pageSize)
}

// CreateDomain 创建域名绑定
func CreateDomain(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req struct {
		InstanceID   uint   `json:"instanceId" binding:"required"`
		DomainName   string `json:"domain" binding:"required"`
		Protocol     string `json:"protocol"`
		InternalIP   string `json:"internalIp" binding:"required"`
		InternalPort int    `json:"internalPort" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "参数错误"))
		return
	}

	if req.Protocol == "" {
		req.Protocol = "http"
	}

	svc := getDomainService()
	d, err := svc.CreateDomain(userID, req.InstanceID, req.DomainName, req.InternalIP, req.InternalPort, req.Protocol, nil)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	common.ResponseSuccess(c, d)
}

// UpdateDomain 更新域名绑定
func UpdateDomain(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	domainID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "域名ID格式错误"))
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "参数错误"))
		return
	}

	svc := getDomainService()
	if err := svc.UpdateDomain(uint(domainID), userID, req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil)
}

// DeleteDomain 删除域名绑定
func DeleteDomain(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	domainID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "域名ID格式错误"))
		return
	}

	svc := getDomainService()
	if err := svc.DeleteDomain(uint(domainID), userID); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil)
}

// GetAvailableQuota 获取可用域名配额
func GetAvailableQuota(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	svc := getDomainService()
	used, max, err := svc.GetAvailableQuota(userID, nil)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "获取配额失败"))
		return
	}

	common.ResponseSuccess(c, gin.H{
		"used":    used,
		"max":     max,
		"remain":  max - int(used),
	})
}
