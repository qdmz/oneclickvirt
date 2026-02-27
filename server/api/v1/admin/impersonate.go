package admin

import (
	"oneclickvirt/global"
	"oneclickvirt/middleware"
	"oneclickvirt/model/admin"
	"oneclickvirt/model/common"
	userModel "oneclickvirt/model/user"
	"oneclickvirt/service/admin/user"
	"oneclickvirt/service/auth"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ImpersonateUser 管理员代用户登录
// @Summary 管理员代用户登录
// @Description 管理员使用指定用户名义登录到用户管理菜单
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} common.Response{data=string} "登录成功，返回新的JWT Token"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 403 {object} common.Response "权限不足"
// @Failure 404 {object} common.Response "用户不存在"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /admin/users/{id}/impersonate [post]
func ImpersonateUser(c *gin.Context) {
	// 检查是否为管理员
	if !requireAdminOnly(c) {
		return
	}

	// 获取目标用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInvalidParam, "无效的用户ID"))
		return
	}

	// 获取当前管理员ID
	currentAdminID, err := getUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, "未找到管理员信息"))
		return
	}

	// 验证目标用户是否存在且不是管理员
	userService := user.NewService()
	targetUser, err := userService.GetUserByID(uint(userID))
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUserNotFound, "用户不存在"))
		return
	}

	// 禁止管理员代管理员登录
	if targetUser.UserType == "admin" {
		common.ResponseWithError(c, common.NewError(common.CodeForbidden, "不能代管理员用户登录"))
		return
	}

	// 记录操作日志
	global.APP_LOG.Info("管理员代用户登录",
		zap.Uint("admin_id", currentAdminID),
		zap.Uint("target_user_id", uint(userID)),
		zap.String("target_username", targetUser.Username))

	// 生成用户登录Token
	jwtService := auth.GetJWTService()
	token, err := jwtService.GenerateUserToken(targetUser.ID, targetUser.Username, targetUser.UserType)
	if err != nil {
		global.APP_LOG.Error("生成用户Token失败",
			zap.Uint("user_id", targetUser.ID),
			zap.Error(err))
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "生成登录凭证失败"))
		return
	}

	// 返回Token
	response := admin.ImpersonateResponse{
		Token:     token,
		ExpiresIn: int(jwtService.GetTokenExpireDuration().Seconds()),
		UserInfo:  targetUser,
	}

	common.ResponseSuccess(c, response, "代登录成功")
}

// GetUserByID 通过ID获取用户信息
func (s *Service) GetUserByID(userID uint) (*userModel.User, error) {
	var user userModel.User
	if err := global.APP_DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}