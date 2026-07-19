package task

import (
	"context"
	"encoding/json"
	"fmt"

	"oneclickvirt/global"
	adminModel "oneclickvirt/model/admin"
	providerModel "oneclickvirt/model/provider"
	productModel "oneclickvirt/model/product"
	userModel "oneclickvirt/model/user"
	userService "oneclickvirt/service/user"

	"go.uber.org/zap"
)

// AutoProvisionVMRequest 自动创建VM请求
type AutoProvisionVMRequest struct {
	UserID       uint                   `json:"userId"`
	ProductID    uint                   `json:"productId"`
	ProductName  string                 `json:"productName"`
	ProductCPU   int                    `json:"productCpu"`
	ProductMemory int                   `json:"productMemory"`
	ProductDisk  int                    `json:"productDisk"`
	ProductBandwidth int                `json:"productBandwidth"`
	VMType       string                 `json:"vmType"`
	OS           string                 `json:"os"`
	OrderID      uint                   `json:"orderId"`
}

// ExecuteAutoProvisionVMTask 执行自动创建VM任务
func (s *TaskService) ExecuteAutoProvisionVMTask(ctx context.Context, task *adminModel.Task) error {
	s.updateTaskProgress(task.ID, 5, "解析任务数据...")

	// 解析任务数据
	var req AutoProvisionVMRequest
	if err := json.Unmarshal([]byte(task.TaskData), &req); err != nil {
		global.APP_LOG.Error("解析自动创建VM任务数据失败",
			zap.Uint("taskId", task.ID),
			zap.Error(err))
		return fmt.Errorf("解析任务数据失败: %v", err)
	}

	global.APP_LOG.Info("自动创建VM任务数据解析成功",
		zap.Uint("taskId", task.ID),
		zap.Uint("userId", req.UserID),
		zap.Uint("productId", req.ProductID),
		zap.String("vmType", req.VMType),
		zap.String("os", req.OS))

	s.updateTaskProgress(task.ID, 10, "获取产品配置...")

	// 获取产品配置
	var product productModel.Product
	if err := global.APP_DB.First(&product, req.ProductID).Error; err != nil {
		global.APP_LOG.Error("获取产品配置失败",
			zap.Uint("productId", req.ProductID),
			zap.Error(err))
		return fmt.Errorf("获取产品配置失败: %v", err)
	}

	s.updateTaskProgress(task.ID, 20, "获取可用Provider...")

	// 获取可用的Provider列表
	userServiceInstance := userService.NewService()
	providers, err := userServiceInstance.GetAvailableProviders(req.UserID)
	if err != nil || len(providers) == 0 {
		global.APP_LOG.Error("获取可用Provider失败",
			zap.Uint("userId", req.UserID),
			zap.Error(err))
		return fmt.Errorf("暂无可用节点，无法自动创建VM")
	}

	// 根据VM类型选择Provider
	var selectedProvider userModel.AvailableProviderResponse
	for _, p := range providers {
		if req.VMType == "" || p.InstanceType == req.VMType {
			selectedProvider = p
			break
		}
	}
	if selectedProvider.ID == 0 {
		// 如果没有匹配的，使用第一个
		selectedProvider = providers[0]
	}

	global.APP_LOG.Info("选择Provider",
		zap.Uint("providerId", selectedProvider.ID),
		zap.String("providerName", selectedProvider.Name))

	s.updateTaskProgress(task.ID, 30, "获取系统镜像...")

	// 获取系统镜像列表
	var systemImages []userModel.SystemImageResponse
	if req.OS != "" {
		systemImages, err = userServiceInstance.GetSystemImages(req.UserID, userModel.SystemImagesRequest{})
		if err != nil {
			systemImages, _ = userServiceInstance.GetFilteredSystemImages(req.UserID, selectedProvider.ID, req.VMType)
		}
	}

	// 选择第一个可用的镜像
	var imageID uint
	if len(systemImages) > 0 {
		imageID = systemImages[0].ID
		global.APP_LOG.Info("选择系统镜像",
			zap.Uint("imageId", imageID),
			zap.String("imageName", systemImages[0].Name))
	}

	s.updateTaskProgress(task.ID, 40, "获取配置选项...")

	// 获取实例配置选项
	instanceConfig, err := userServiceInstance.GetInstanceConfig(req.UserID, selectedProvider.ID)
	if err != nil {
		global.APP_LOG.Error("获取实例配置选项失败",
			zap.Error(err))
		return fmt.Errorf("获取配置选项失败: %v", err)
	}

	// 从产品配置中选择规格
	cpuSpec := findClosestSpec(instanceConfig.CPUs, req.ProductCPU)
	memorySpec := findClosestSpec(instanceConfig.Memories, req.ProductMemory)
	diskSpec := findClosestSpec(instanceConfig.Disks, req.ProductDisk)
	bandwidthSpec := findClosestSpec(instanceConfig.Bandwidths, req.ProductBandwidth)

	global.APP_LOG.Info("选择的规格",
		zap.String("cpu", cpuSpec.ID),
		zap.String("memory", memorySpec.ID),
		zap.String("disk", diskSpec.ID),
		zap.String("bandwidth", bandwidthSpec.ID))

	s.updateTaskProgress(task.ID, 50, "创建VM实例...")

	// 构建创建请求
	createReq := userModel.CreateInstanceRequest{
		ProviderId:  selectedProvider.ID,
		ImageId:     imageID,
		CPUId:       cpuSpec.ID,
		MemoryId:    memorySpec.ID,
		DiskId:      diskSpec.ID,
		BandwidthId: bandwidthSpec.ID,
		Description: fmt.Sprintf("自动创建 - 产品: %s", product.Name),
	}

	// 创建VM实例
	taskResult, err := userServiceInstance.CreateUserInstance(req.UserID, createReq)
	if err != nil {
		global.APP_LOG.Error("创建VM实例失败",
			zap.Uint("userId", req.UserID),
			zap.Error(err))
		return fmt.Errorf("创建VM失败: %v", err)
	}

	global.APP_LOG.Info("VM实例创建任务已提交",
		zap.Uint("taskId", taskResult.ID),
		zap.Uint("userId", req.UserID))

	s.updateTaskProgress(task.ID, 90, "VM创建任务已提交，等待完成...")

	// 更新订单状态
	global.APP_LOG.Info("更新订单状态为已处理",
		zap.Uint("orderId", req.OrderID))

	var order adminModel.Order
	if err := global.APP_DB.First(&order, req.OrderID).Error; err == nil {
		order.AutoProvisionTaskID = &taskResult.ID
		global.APP_DB.Save(&order)
	}

	s.updateTaskProgress(task.ID, 100, "自动创建VM任务完成")
	return nil
}

// findClosestSpec 从规格列表中找到最接近的规格
func findClosestSpec(specs []userModel.InstanceConfigSpec, targetValue int) userModel.InstanceConfigSpec {
	if len(specs) == 0 {
		return userModel.InstanceConfigSpec{}
	}

	bestMatch := specs[0]
	bestDiff := abs(bestMatch.Value - targetValue)

	for _, spec := range specs {
		diff := abs(spec.Value - targetValue)
		if diff < bestDiff {
			bestDiff = diff
			bestMatch = spec
		}
	}

	return bestMatch
}

// abs 绝对值函数
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
