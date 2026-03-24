package admin

import (
	"fmt"
	"oneclickvirt/global"
	redemptionModel "oneclickvirt/model/redemption"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateRedemptionCodeFixed 兑换码创建（修复版）
// 修复了金额单位转换和验证问题
func CreateRedemptionCodeFixed(c *gin.Context) {
	var req struct {
		Code      string     `json:"code"`
		Type      string     `json:"type" binding:"required"`
		Amount    int64      `json:"amount"`         // 前端已经将余额类型转换为分
		ProductID *uint      `json:"productId"`
		MaxUses   int        `json:"maxUses"`
		ExpireAt  *time.Time `json:"expireAt"`
		Remark    string     `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		global.APP_LOG.Error("解析兑换码参数失败", zap.Error(err))
		c.JSON(400, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 生成兑换码(如果为空)
	if req.Code == "" {
		req.Code = generateRedemptionCode()
	}

	// 验证必填字段
	if req.Type == "" {
		c.JSON(400, gin.H{"code": 400, "message": "兑换码类型不能为空"})
		return
	}

	// 如果是产品类型，必须有产品ID
	if req.Type == redemptionModel.RedemptionTypeProduct && req.ProductID == nil {
		c.JSON(400, gin.H{"code": 400, "message": "产品类型必须指定产品ID"})
		return
	}

	// 处理金额
	var amount int64
	if req.Type == redemptionModel.RedemptionTypeBalance {
		// 余额类型：前端已经将金额转换为分
		amount = req.Amount
	} else if req.Type == redemptionModel.RedemptionTypeLevel {
		// 等级类型：直接使用整数
		amount = req.Amount
	} else if req.Type == redemptionModel.RedemptionTypeProduct {
		// 产品类型：金额为0
		amount = 0
	}

	// 创建兑换码
	code := redemptionModel.RedemptionCode{
		Code:      req.Code,
		Type:      req.Type,
		Amount:    amount,
		ProductID: req.ProductID,
		MaxUses:   req.MaxUses,
		ExpireAt:  req.ExpireAt,
		Remark:    req.Remark,
		IsEnabled: true,
	}

	// 设置默认最大使用次数
	if code.MaxUses == 0 {
		code.MaxUses = 1
	}

	global.APP_LOG.Info("创建兑换码", zap.String("code", code.Code), zap.String("type", code.Type), zap.Int64("amount", code.Amount))

	if err := global.APP_DB.Create(&code).Error; err != nil {
		global.APP_LOG.Error("创建兑换码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "创建兑换码失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    code,
	})
}

// GenerateRedemptionCodesFixed 批量生成兑换码（修复版）
func GenerateRedemptionCodesFixed(c *gin.Context) {
	var params struct {
		Count    int        `json:"count" binding:"required"`
		Type     string     `json:"type" binding:"required"`
		Amount   int64      `json:"amount"`         // 前端已经将余额类型转换为分
		MaxUses  int        `json:"maxUses"`
		ExpireAt *time.Time `json:"expireAt"`
		Remark   string     `json:"remark"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		global.APP_LOG.Error("解析批量生成参数失败", zap.Error(err))
		c.JSON(400, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if params.Count <= 0 || params.Count > 1000 {
		c.JSON(400, gin.H{"code": 400, "message": "生成数量必须在1-1000之间"})
		return
	}

	// 处理金额
	var amount int64
	if params.Type == redemptionModel.RedemptionTypeBalance {
		// 余额类型：前端已经将金额转换为分
		amount = params.Amount
	} else if params.Type == redemptionModel.RedemptionTypeLevel {
		// 等级类型：直接使用整数
		amount = params.Amount
	} else if params.Type == redemptionModel.RedemptionTypeProduct {
		// 产品类型：金额为0
		amount = 0
	}

	// 生成兑换码
	codes := make([]redemptionModel.RedemptionCode, params.Count)
	for i := 0; i < params.Count; i++ {
		codes[i] = redemptionModel.RedemptionCode{
			Code:      generateRedemptionCode(),
			Type:      params.Type,
			Amount:    amount,
			MaxUses:   params.MaxUses,
			ExpireAt:  params.ExpireAt,
			Remark:    params.Remark,
			IsEnabled: true,
		}
		if params.MaxUses == 0 {
			codes[i].MaxUses = 1
		}
	}

	global.APP_LOG.Info("开始批量创建兑换码", zap.Int("count", params.Count), zap.Int64("amount", amount))

	if err := global.APP_DB.CreateInBatches(&codes, 100).Error; err != nil {
		global.APP_LOG.Error("批量生成兑换码失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "批量生成兑换码失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": fmt.Sprintf("成功生成%d个兑换码", params.Count),
		"data":    codes,
	})
}
