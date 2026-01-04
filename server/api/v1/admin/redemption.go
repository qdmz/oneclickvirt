package admin

import (
	"fmt"
	"oneclickvirt/global"
	redemptionModel "oneclickvirt/model/redemption"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetRedemptionCodes 获取兑换码列表
// @Summary 获取兑换码列表
// @Description 管理员获取兑换码列表
// @Tags 管理员/兑换码管理
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /v1/admin/redemption-codes [get]
func GetRedemptionCodes(c *gin.Context) {
	var codes []redemptionModel.RedemptionCode
	query := global.APP_DB.Order("created_at DESC")

	// 搜索条件
	if code := c.Query("code"); code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if typeStr := c.Query("type"); typeStr != "" {
		query = query.Where("type = ?", typeStr)
	}
	if enabled := c.Query("enabled"); enabled != "" {
		if enabled == "true" {
			query = query.Where("is_enabled = ?", true)
		} else {
			query = query.Where("is_enabled = ?", false)
		}
	}

	if err := query.Find(&codes).Error; err != nil {
		global.APP_LOG.Error("获取兑换码列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取兑换码列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": codes,
	})
}

// CreateRedemptionCode 创建兑换码
// @Summary 创建兑换码
// @Description 管理员创建单个兑换码
// @Tags 管理员/兑换码管理
// @Accept json
// @Produce json
// @Param code body redemptionModel.RedemptionCode true "兑换码信息"
// @Success 200 {object} common.Response
// @Router /v1/admin/redemption-codes [post]
func CreateRedemptionCode(c *gin.Context) {
	var code redemptionModel.RedemptionCode
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 生成兑换码(如果为空)
	if code.Code == "" {
		code.Code = generateRedemptionCode()
	}

	// 验证必填字段
	if code.Type == "" {
		c.JSON(400, gin.H{"code": 400, "message": "兑换码类型不能为空"})
		return
	}

	if err := global.APP_DB.Create(&code).Error; err != nil {
		global.APP_LOG.Error("创建兑换码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "创建兑换码失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "创建成功",
		"data": code,
	})
}

// GenerateRedemptionCodes 批量生成兑换码
// @Summary 批量生成兑换码
// @Description 管理员批量生成兑换码
// @Tags 管理员/兑换码管理
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "生成参数: count(数量), type(类型), amount(金额), maxUses(最大使用次数), expireAt(过期时间)"
// @Success 200 {object} common.Response
// @Router /v1/admin/redemption-codes/generate [post]
func GenerateRedemptionCodes(c *gin.Context) {
	var params struct {
		Count    int        `json:"count" binding:"required"`
		Type     string     `json:"type" binding:"required"`
		Amount   int64      `json:"amount"`
		MaxUses  int        `json:"maxUses"`
		ExpireAt *time.Time `json:"expireAt"`
		Remark   string     `json:"remark"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if params.Count <= 0 || params.Count > 1000 {
		c.JSON(400, gin.H{"code": 400, "message": "生成数量必须在1-1000之间"})
		return
	}

	// 生成兑换码
	codes := make([]redemptionModel.RedemptionCode, params.Count)
	for i := 0; i < params.Count; i++ {
		codes[i] = redemptionModel.RedemptionCode{
			Code:     generateRedemptionCode(),
			Type:     params.Type,
			Amount:   params.Amount,
			MaxUses:  params.MaxUses,
			ExpireAt: params.ExpireAt,
			Remark:   params.Remark,
			IsEnabled: true,
		}
		if params.MaxUses == 0 {
			codes[i].MaxUses = 1
		}
	}

	if err := global.APP_DB.CreateInBatches(&codes, 100).Error; err != nil {
		global.APP_LOG.Error("批量生成兑换码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "批量生成兑换码失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": fmt.Sprintf("成功生成%d个兑换码", params.Count),
		"data": codes,
	})
}

// DeleteRedemptionCode 删除兑换码
// @Summary 删除兑换码
// @Description 管理员删除兑换码
// @Tags 管理员/兑换码管理
// @Accept json
// @Produce json
// @Param id path uint true "兑换码ID"
// @Success 200 {object} common.Response
// @Router /v1/admin/redemption-codes/{id} [delete]
func DeleteRedemptionCode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "兑换码ID不能为空"})
		return
	}

	if err := global.APP_DB.Delete(&redemptionModel.RedemptionCode{}, id).Error; err != nil {
		global.APP_LOG.Error("删除兑换码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "删除兑换码失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "删除成功",
	})
}

// GetRedemptionCodeUsages 获取兑换码使用记录
// @Summary 获取兑换码使用记录
// @Description 管理员获取兑换码使用记录
// @Tags 管理员/兑换码管理
// @Accept json
// @Produce json
// @Param id path uint true "兑换码ID"
// @Success 200 {object} common.Response
// @Router /v1/admin/redemption-codes/{id}/usages [get]
func GetRedemptionCodeUsages(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "兑换码ID不能为空"})
		return
	}

	var usages []redemptionModel.RedemptionCodeUsage
	if err := global.APP_DB.Where("code_id = ?", id).
		Preload("User").
		Order("created_at DESC").
		Find(&usages).Error; err != nil {
		global.APP_LOG.Error("获取兑换码使用记录失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取使用记录失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "success",
		"data": usages,
	})
}

// ToggleRedemptionCode 启用/禁用兑换码
// @Summary 启用/禁用兑换码
// @Description 管理员启用或禁用兑换码
// @Tags 管理员/兑换码管理
// @Accept json
// @Produce json
// @Param id path uint true "兑换码ID"
// @Success 200 {object} common.Response
// @Router /v1/admin/redemption-codes/{id}/toggle [put]
func ToggleRedemptionCode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"code": 400, "message": "兑换码ID不能为空"})
		return
	}

	var code redemptionModel.RedemptionCode
	if err := global.APP_DB.First(&code, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"code": 404, "message": "兑换码不存在"})
			return
		}
		global.APP_LOG.Error("查询兑换码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "查询兑换码失败"})
		return
	}

	code.IsEnabled = !code.IsEnabled
	if err := global.APP_DB.Save(&code).Error; err != nil {
		global.APP_LOG.Error("更新兑换码状态失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "更新兑换码状态失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"message": "操作成功",
		"data": code,
	})
}

// generateRedemptionCode 生成兑换码
func generateRedemptionCode() string {
	// 生成16位随机兑换码
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(time.Nanosecond)
	}
	return string(b[0:4]) + "-" + string(b[4:8]) + "-" + string(b[8:12]) + "-" + string(b[12:16])
}
