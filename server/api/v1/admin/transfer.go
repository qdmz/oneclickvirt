package admin

import (
	"oneclickvirt/global"
	"oneclickvirt/model/admin"
	"oneclickvirt/model/common"
	userModel "oneclickvirt/model/user"
	providerModel "oneclickvirt/model/provider"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TransferInstanceOwnership 转移实例归属
// @Summary 转移实例归属
// @Description 管理员将实例转移给其他用户
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "实例ID"
// @Param request body admin.TransferInstanceRequest true "转移请求参数"
// @Success 200 {object} common.Response "转移成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 403 {object} common.Response "权限不足"
// @Failure 404 {object} common.Response "实例或用户不存在"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /admin/instances/{id}/transfer [post]
func TransferInstanceOwnership(c *gin.Context) {
	// 检查是否为管理员
	if !requireAdminOnly(c) {
		return
	}

	// 获取实例ID
	instanceIDStr := c.Param("id")
	instanceID, err := strconv.ParseUint(instanceIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "无效的实例ID"))
		return
	}

	// 解析请求参数
	var req admin.TransferInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	// 验证目标用户是否存在
	var targetUser userModel.User
	if err := global.APP_DB.First(&targetUser, req.TargetUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			common.ResponseWithError(c, common.NewError(common.CodeUserNotFound, "目标用户不存在"))
			return
		}
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "查询用户失败"))
		return
	}

	// 禁止将实例转移给管理员用户
	if targetUser.UserType == "admin" {
		common.ResponseWithError(c, common.NewError(common.CodeForbidden, "不能将实例转移给管理员用户"))
		return
	}

	// 获取实例信息
	var instance providerModel.Instance
	if err := global.APP_DB.First(&instance, uint(instanceID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			common.ResponseWithError(c, common.NewError(common.CodeNotFound, "实例不存在"))
			return
		}
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "查询实例失败"))
		return
	}

	// 检查实例是否属于当前用户（双重检查）
	if instance.UserID == targetUser.ID {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "实例已属于目标用户"))
		return
	}

	// 记录操作日志
	global.APP_LOG.Info("管理员转移实例归属",
		zap.Uint("instance_id", uint(instanceID)),
		zap.Uint("from_user_id", instance.UserID),
		zap.Uint("to_user_id", targetUser.ID),
		zap.String("instance_name", instance.Name))

	// 执行转移操作
	if err := global.APP_DB.Model(&instance).Update("user_id", targetUser.ID).Error; err != nil {
		global.APP_LOG.Error("转移实例归属失败",
			zap.Uint("instance_id", uint(instanceID)),
			zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "转移失败"))
		return
	}

	// 如果实例有到期时间，更新为目标用户的有效期
	if targetUser.LevelExpireAt != nil {
		if err := global.APP_DB.Model(&instance).Update("expired_at", targetUser.LevelExpireAt).Error; err != nil {
			global.APP_LOG.Warn("更新实例到期时间失败",
				zap.Uint("instance_id", uint(instanceID)),
				zap.Error(err))
		}
	}

	common.ResponseSuccess(c, nil, "实例转移成功")
}