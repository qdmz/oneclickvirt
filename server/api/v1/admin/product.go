package admin

import (
	"fmt"
	"oneclickvirt/global"
	productModel "oneclickvirt/model/product"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetProducts 获取产品列表
// @Summary 获取产品列表
// @Description 管理员获取产品列表
// @Tags 管理员/产品管理
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /v1/admin/products [get]
func GetProducts(c *gin.Context) {
	var products []productModel.Product
	query := global.APP_DB.Order("sort_order ASC, id ASC")

	// 搜索条件
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if level := c.Query("level"); level != "" {
		query = query.Where("level = ?", level)
	}
	if enabled := c.Query("enabled"); enabled != "" {
		if enabled == "true" {
			query = query.Where("is_enabled = ?", 1)
		} else {
			query = query.Where("is_enabled = ?", 0)
		}
	}

	if err := query.Find(&products).Error; err != nil {
		global.APP_LOG.Error("获取产品列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取产品列表失败"})
		return
	}

	// 从用户等级限制中自动填充产品资源配置，确保产品配置与最新的等级配额一致
	for i := range products {
		fillProductResourcesFromLevelLimit(&products[i])
	}

	global.APP_LOG.Info("获取产品列表成功", zap.Int("count", len(products)))
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    products,
	})
}

// CreateProduct 创建产品
// @Summary 创建产品
// @Description 管理员创建新产品
// @Tags 管理员/产品管理
// @Accept json
// @Produce json
// @Param product body productModel.Product true "产品信息"
// @Success 200 {object} common.Response
// @Router /v1/admin/products [post]
func CreateProduct(c *gin.Context) {
	var product productModel.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证必填字段
	if product.Name == "" || product.Level <= 0 {
		c.JSON(400, gin.H{"code": 400, "message": "产品名称和等级不能为空"})
		return
	}

	// 从用户等级限制中自动填充产品资源配置
	fillProductResourcesFromLevelLimit(&product)

	if err := global.APP_DB.Create(&product).Error; err != nil {
		global.APP_LOG.Error("创建产品失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "创建产品失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    product,
	})
}

// UpdateProduct 更新产品
// @Summary 更新产品
// @Description 管理员更新产品信息
// @Tags 管理员/产品管理
// @Accept json
// @Produce json
// @Param id path uint true "产品ID"
// @Param product body productModel.Product true "产品信息"
// @Success 200 {object} common.Response
// @Router /v1/admin/products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "产品ID不能为空"})
		return
	}

	var product productModel.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 从用户等级限制中自动填充产品资源配置
	fillProductResourcesFromLevelLimit(&product)

	// 创建一个map来存储需要更新的字段，确保is_enabled字段被正确处理
	updateMap := map[string]interface{}{
		"name":          product.Name,
		"description":   product.Description,
		"level":         product.Level,
		"price":         product.Price,
		"period":        product.Period,
		"cpu":           product.CPU,
		"memory":        product.Memory,
		"disk":          product.Disk,
		"bandwidth":     product.Bandwidth,
		"traffic":       product.Traffic,
		"max_instances": product.MaxInstances,
		"is_enabled":    product.IsEnabled, // 明确指定is_enabled字段，确保0值也能被更新
		"sort_order":    product.SortOrder,
		"features":      product.Features,
		"allow_repeat":  product.AllowRepeat, // 添加是否允许重复购买字段
	}

	// 更新产品
	if err := global.APP_DB.Model(&productModel.Product{}).Where("id = ?", id).Updates(updateMap).Error; err != nil {
		global.APP_LOG.Error("更新产品失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "更新产品失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteProduct 删除产品
// @Summary 删除产品
// @Description 管理员删除产品
// @Tags 管理员/产品管理
// @Accept json
// @Produce json
// @Param id path uint true "产品ID"
// @Success 200 {object} common.Response
// @Router /v1/admin/products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "产品ID不能为空"})
		return
	}

	if err := global.APP_DB.Delete(&productModel.Product{}, id).Error; err != nil {
		global.APP_LOG.Error("删除产品失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "删除产品失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ToggleProduct 启用/禁用产品
// @Summary 启用/禁用产品
// @Description 管理员启用或禁用产品
// @Tags 管理员/产品管理
// @Accept json
// @Produce json
// @Param id path uint true "产品ID"
// @Success 200 {object} common.Response
// @Router /v1/admin/products/{id}/toggle [put]
func ToggleProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "产品ID不能为空"})
		return
	}

	var product productModel.Product
	if err := global.APP_DB.First(&product, id).Error; err != nil {
		global.APP_LOG.Error("产品不存在", zap.Error(err))
		c.JSON(404, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	// 切换产品状态（1:启用, 0:禁用）
	newStatus := 0
	if product.IsEnabled == 0 {
		newStatus = 1
	}

	// 使用Updates明确更新is_enabled字段，确保0值也能被更新
	if err := global.APP_DB.Model(&product).Updates(map[string]interface{}{
		"is_enabled": newStatus,
	}).Error; err != nil {
		global.APP_LOG.Error("更新产品状态失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "更新产品状态失败"})
		return
	}

	// 重新获取更新后的产品数据
	if err := global.APP_DB.First(&product, id).Error; err != nil {
		global.APP_LOG.Error("获取更新后的产品数据失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "操作成功，但获取更新数据失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    product,
	})
}

// fillProductResourcesFromLevelLimit 从用户等级限制中自动填充产品资源配置
func fillProductResourcesFromLevelLimit(product *productModel.Product) {
	// 获取全局配置中的等级限制
	levelLimits := global.APP_CONFIG.Quota.LevelLimits
	level := product.Level

	if levelInfo, exists := levelLimits[level]; exists {
		// 填充最大实例数
		product.MaxInstances = levelInfo.MaxInstances

		// 填充最大流量（直接使用int64类型，无需转换）
		product.Traffic = levelInfo.MaxTraffic

		// 填充资源配置
		maxResources := levelInfo.MaxResources
		if maxResources != nil {
			// 填充CPU核心数
			if cpu, ok := maxResources["cpu"].(float64); ok {
				product.CPU = int(cpu)
			} else if cpu, ok := maxResources["cpu"].(int); ok {
				product.CPU = cpu
			}

			// 填充内存(MB)
			if memory, ok := maxResources["memory"].(float64); ok {
				product.Memory = int(memory)
			} else if memory, ok := maxResources["memory"].(int); ok {
				product.Memory = memory
			}

			// 填充磁盘(MB)
			if disk, ok := maxResources["disk"].(float64); ok {
				product.Disk = int(disk)
			} else if disk, ok := maxResources["disk"].(int); ok {
				product.Disk = disk
			}

			// 填充带宽(Mbps)
			if bandwidth, ok := maxResources["bandwidth"].(float64); ok {
				product.Bandwidth = int(bandwidth)
			} else if bandwidth, ok := maxResources["bandwidth"].(int); ok {
				product.Bandwidth = bandwidth
			}
		}

		global.APP_LOG.Info(fmt.Sprintf("产品 %s (等级 %d) 资源配置已从等级限制自动填充", product.Name, level))
	} else {
		global.APP_LOG.Warn(fmt.Sprintf("未找到等级 %d 的限制配置，无法自动填充产品资源配置", level))
		// 记录可用的等级限制键
		availableLevels := make([]int, 0, len(levelLimits))
		for k := range levelLimits {
			availableLevels = append(availableLevels, k)
		}
		global.APP_LOG.Warn(fmt.Sprintf("可用的等级限制键: %v", availableLevels))
	}
}
