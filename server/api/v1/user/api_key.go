package user

import (
	"net/http"
	"oneclickvirt/global"
	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	userModel "oneclickvirt/model/user"
	"oneclickvirt/service/auth"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// APIKeyController API 密钥控制器
type APIKeyController struct{}

// CreateAPIKey 创建 API 密钥
// @Summary 创建 API 密钥
// @Description 为当前用户创建一个新的 API 密钥
// @Tags API 密钥管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.CreateAPIKeyRequest true "创建请求"
// @Success 200 {object} common.Response{data=user.APIKey} "创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 500 {object} common.Response "创建失败"
// @Router /api/v1/user/api-keys [post]
func (c *APIKeyController) CreateAPIKey(ctx *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	var req auth.CreateAPIKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}

	service := auth.NewAPIKeyService()
	apiKey, err := service.CreateAPIKey(authCtx.UserID, req)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(ctx, apiKey, "创建成功")
}

// GetAPIKeys 获取 API 密钥列表
// @Summary 获取 API 密钥列表
// @Description 获取当前用户的所有 API 密钥
// @Tags API 密钥管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response{data=[]user.APIKey} "获取成功"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 500 {object} common.Response "获取失败"
// @Router /api/v1/user/api-keys [get]
func (c *APIKeyController) GetAPIKeys(ctx *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	service := auth.NewAPIKeyService()
	apiKeys, err := service.GetAPIKeys(authCtx.UserID)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(ctx, apiKeys)
}

// GetAPIKey 获取 API 密钥详情
// @Summary 获取 API 密钥详情
// @Description 获取指定 ID 的 API 密钥详情
// @Tags API 密钥管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "密钥 ID"
// @Success 200 {object} common.Response{data=user.APIKey} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 404 {object} common.Response "密钥不存在"
// @Failure 500 {object} common.Response "获取失败"
// @Router /api/v1/user/api-keys/{id} [get]
func (c *APIKeyController) GetAPIKey(ctx *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	idUint := uint(id)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := auth.NewAPIKeyService()
	apiKey, err := service.GetAPIKeyByID(idUint, authCtx.UserID)
	if err != nil {
		if err.Error() == "密钥不存在" {
			ctx.JSON(http.StatusNotFound, common.Response{
				Code: 404,
				Msg:  "密钥不存在",
			})
			return
		}
		common.ResponseWithError(ctx, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(ctx, apiKey)
}

// UpdateAPIKey 更新 API 密钥
// @Summary 更新 API 密钥
// @Description 更新指定 ID 的 API 密钥
// @Tags API 密钥管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "密钥 ID"
// @Param request body auth.UpdateAPIKeyRequest true "更新请求"
// @Success 200 {object} common.Response{data=user.APIKey} "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 404 {object} common.Response "密钥不存在"
// @Failure 500 {object} common.Response "更新失败"
// @Router /api/v1/user/api-keys/{id} [put]
func (c *APIKeyController) UpdateAPIKey(ctx *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	idUint := uint(id)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	var req auth.UpdateAPIKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}

	service := auth.NewAPIKeyService()
	apiKey, err := service.UpdateAPIKey(idUint, authCtx.UserID, req)
	if err != nil {
		if err.Error() == "密钥不存在" {
			ctx.JSON(http.StatusNotFound, common.Response{
				Code: 404,
				Msg:  "密钥不存在",
			})
			return
		}
		common.ResponseWithError(ctx, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(ctx, apiKey, "更新成功")
}

// DeleteAPIKey 删除 API 密钥
// @Summary 删除 API 密钥
// @Description 删除指定 ID 的 API 密钥
// @Tags API 密钥管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "密钥 ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 404 {object} common.Response "密钥不存在"
// @Failure 500 {object} common.Response "删除失败"
// @Router /api/v1/user/api-keys/{id} [delete]
func (c *APIKeyController) DeleteAPIKey(ctx *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	idUint := uint(id)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := auth.NewAPIKeyService()
	if err := service.DeleteAPIKey(idUint, authCtx.UserID); err != nil {
		if err.Error() == "密钥不存在" {
			ctx.JSON(http.StatusNotFound, common.Response{
				Code: 404,
				Msg:  "密钥不存在",
			})
			return
		}
		common.ResponseWithError(ctx, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(ctx, nil, "删除成功")
}

// RevokeAPIKey 撤销 API 密钥
// @Summary 撤销 API 密钥
// @Description 撤销（禁用）指定 ID 的 API 密钥
// @Tags API 密钥管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "密钥 ID"
// @Success 200 {object} common.Response "撤销成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 404 {object} common.Response "密钥不存在"
// @Failure 500 {object} common.Response "撤销失败"
// @Router /api/v1/user/api-keys/{id}/revoke [put]
func (c *APIKeyController) RevokeAPIKey(ctx *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	idUint := uint(id)
	if err != nil {
		common.ResponseWithError(ctx, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	service := auth.NewAPIKeyService()
	if err := service.RevokeAPIKey(idUint, authCtx.UserID); err != nil {
		if err.Error() == "密钥不存在" {
			ctx.JSON(http.StatusNotFound, common.Response{
				Code: 404,
				Msg:  "密钥不存在",
			})
			return
		}
		common.ResponseWithError(ctx, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(ctx, nil, "撤销成功")
}

// GetAPIKeyStats 获取 API 密钥统计信息
// @Summary 获取 API 密钥统计信息
// @Description 获取当前用户的 API 密钥统计信息
// @Tags API 密钥管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.Response{data=map[string]interface{}} "获取成功"
// @Failure 401 {object} common.Response "认证失败"
// @Failure 500 {object} common.Response "获取失败"
// @Router /api/v1/user/api-keys/stats [get]
func (c *APIKeyController) GetAPIKeyStats(ctx *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.Response{
			Code: 401,
			Msg:  "用户未认证",
		})
		return
	}

	var stats struct {
		Total      int64     `json:"total"`
		Active     int64     `json:"active"`
		Disabled   int64     `json:"disabled"`
		Expired    int64     `json:"expired"`
		LastUsedAt *time.Time `json:"lastUsedAt"`
	}

	// 获取总数
	global.APP_DB.Model(&userModel.APIKey{}).Where("user_id = ?", authCtx.UserID).Count(&stats.Total)

	// 获取活跃数量
	global.APP_DB.Model(&userModel.APIKey{}).Where("user_id = ? AND status = ?", authCtx.UserID, "active").Count(&stats.Active)

	// 获取禁用数量
	global.APP_DB.Model(&userModel.APIKey{}).Where("user_id = ? AND status = ?", authCtx.UserID, "disabled").Count(&stats.Disabled)

	// 获取过期数量
	global.APP_DB.Model(&userModel.APIKey{}).Where("user_id = ? AND expires_at < ?", authCtx.UserID, time.Now()).Count(&stats.Expired)

	// 获取最后使用时间
	var lastUsed userModel.APIKey
	global.APP_DB.Model(&userModel.APIKey{}).
		Where("user_id = ? AND last_used_at IS NOT NULL", authCtx.UserID).
		Order("last_used_at DESC").
		First(&lastUsed)
	if lastUsed.ID > 0 {
		stats.LastUsedAt = lastUsed.LastUsedAt
	}

	common.ResponseSuccess(ctx, stats)
}
