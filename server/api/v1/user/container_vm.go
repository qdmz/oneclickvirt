package user

import (
	"strconv"

	"oneclickvirt/model/common"
	"oneclickvirt/model/user"
	userService "oneclickvirt/service/user"

	"github.com/gin-gonic/gin"
)

// GetUserContainers 获取用户容器列表
// @Summary 获取用户容器列表
// @Description 获取当前用户的所有容器实例
// @Tags 用户/容器管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query string false "实例状态"
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "获取失败"
// @Router /user/containers [get]
func GetUserContainers(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req user.UserInstanceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	req.InstanceType = "container"

	userServiceInstance := userService.NewService()
	instances, total, err := userServiceInstance.GetUserInstances(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "获取容器列表失败"))
		return
	}

	common.ResponseSuccessWithPagination(c, instances, total, req.Page, req.PageSize)
}

// GetUserVMs 获取用户虚拟机列表
// @Summary 获取用户虚拟机列表
// @Description 获取当前用户的所有虚拟机实例
// @Tags 用户/虚拟机管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query string false "实例状态"
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "获取失败"
// @Router /user/vms [get]
func GetUserVMs(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req user.UserInstanceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	req.InstanceType = "vm"

	userServiceInstance := userService.NewService()
	instances, total, err := userServiceInstance.GetUserInstances(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "获取虚拟机列表失败"))
		return
	}

	common.ResponseSuccessWithPagination(c, instances, total, req.Page, req.PageSize)
}

// CreateUserContainer 创建容器
// @Summary 创建容器
// @Description 用户创建新的容器实例（异步处理）
// @Tags 用户/容器管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.CreateInstanceRequest true "创建容器请求参数"
// @Success 200 {object} common.Response{data=object} "任务创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "创建失败"
// @Router /user/containers [post]
func CreateUserContainer(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req user.CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}

	userServiceInstance := userService.NewService()
	task, err := userServiceInstance.CreateUserInstance(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	responseData := map[string]interface{}{
		"taskId":  task.ID,
		"status":  task.Status,
		"message": "容器创建任务已提交，正在后台处理",
	}

	common.ResponseSuccess(c, responseData, "容器创建任务已提交")
}

// ControlUserContainer 控制容器
// @Summary 控制容器
// @Description 对用户容器执行操作（启动、停止、重启等）
// @Tags 用户/容器管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "容器ID"
// @Param request body user.InstanceActionRequest true "容器操作请求参数"
// @Success 200 {object} common.Response "操作成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "操作失败"
// @Router /user/containers/{id}/action [post]
func ControlUserContainer(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	instanceIDStr := c.Param("id")
	instanceID, err := strconv.ParseUint(instanceIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的容器ID"))
		return
	}

	var req user.InstanceActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	req.InstanceID = uint(instanceID)

	userServiceInstance := userService.NewService()
	err = userServiceInstance.InstanceAction(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "操作成功")
}

// DeleteUserContainer 删除容器
// @Summary 删除容器
// @Description 删除用户的容器实例
// @Tags 用户/容器管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "容器ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "删除失败"
// @Router /user/containers/{id} [delete]
func DeleteUserContainer(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	instanceIDStr := c.Param("id")
	instanceID, err := strconv.ParseUint(instanceIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的容器ID"))
		return
	}

	var req user.InstanceActionRequest
	req.InstanceID = uint(instanceID)
	req.Action = "delete"

	userServiceInstance := userService.NewService()
	err = userServiceInstance.InstanceAction(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "删除成功")
}

// CreateUserVM 创建虚拟机
// @Summary 创建虚拟机
// @Description 用户创建新的虚拟机实例（异步处理）
// @Tags 用户/虚拟机管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.CreateInstanceRequest true "创建虚拟机请求参数"
// @Success 200 {object} common.Response{data=object} "任务创建成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "创建失败"
// @Router /user/vms [post]
func CreateUserVM(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req user.CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}

	userServiceInstance := userService.NewService()
	task, err := userServiceInstance.CreateUserInstance(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	responseData := map[string]interface{}{
		"taskId":  task.ID,
		"status":  task.Status,
		"message": "虚拟机创建任务已提交，正在后台处理",
	}

	common.ResponseSuccess(c, responseData, "虚拟机创建任务已提交")
}

// ControlUserVM 控制虚拟机
// @Summary 控制虚拟机
// @Description 对用户虚拟机执行操作（启动、停止、重启等）
// @Tags 用户/虚拟机管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "虚拟机ID"
// @Param request body user.InstanceActionRequest true "虚拟机操作请求参数"
// @Success 200 {object} common.Response "操作成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "操作失败"
// @Router /user/vms/{id}/action [post]
func ControlUserVM(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	instanceIDStr := c.Param("id")
	instanceID, err := strconv.ParseUint(instanceIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的虚拟机ID"))
		return
	}

	var req user.InstanceActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}
	req.InstanceID = uint(instanceID)

	userServiceInstance := userService.NewService()
	err = userServiceInstance.InstanceAction(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "操作成功")
}

// DeleteUserVM 删除虚拟机
// @Summary 删除虚拟机
// @Description 删除用户的虚拟机实例
// @Tags 用户/虚拟机管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "虚拟机ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "删除失败"
// @Router /user/vms/{id} [delete]
func DeleteUserVM(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	instanceIDStr := c.Param("id")
	instanceID, err := strconv.ParseUint(instanceIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的虚拟机ID"))
		return
	}

	var req user.InstanceActionRequest
	req.InstanceID = uint(instanceID)
	req.Action = "delete"

	userServiceInstance := userService.NewService()
	err = userServiceInstance.InstanceAction(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, nil, "删除成功")
}

// GetInstanceLogs 获取实例日志
// @Summary 获取实例日志
// @Description 获取用户实例的控制台日志
// @Tags 用户/实例管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "实例ID"
// @Param lines query int false "日志行数" default(100)
// @Success 200 {object} common.Response{data=object} "获取成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "获取失败"
// @Router /user/instances/{id}/logs [get]
func GetInstanceLogs(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	instanceIDStr := c.Param("id")
	instanceID, err := strconv.ParseUint(instanceIDStr, 10, 32)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的实例ID"))
		return
	}

	lines := 100
	if l := c.Query("lines"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 10000 {
			lines = parsed
		}
	}

	userServiceInstance := userService.NewService()
	logs, err := userServiceInstance.GetInstanceLogs(userID, uint(instanceID), lines)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, err.Error()))
		return
	}

	common.ResponseSuccess(c, logs)
}

// UpdateNickname 更新用户昵称
// @Summary 更新用户昵称
// @Description 更新当前用户的昵称
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "昵称信息: nickname"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "更新失败"
// @Router /user/nickname [put]
func UpdateNickname(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}

	var req user.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误"))
		return
	}

	if req.Nickname == "" {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "昵称不能为空"))
		return
	}

	userServiceInstance := userService.NewService()
	err = userServiceInstance.UpdateProfile(userID, req)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeInternalError, "更新昵称失败"))
		return
	}

	common.ResponseSuccess(c, nil, "更新成功")
}
