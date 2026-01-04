package user

import (
	"fmt"
	"oneclickvirt/global"
	walletModel "oneclickvirt/model/wallet"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetWallet 获取用户钱包信息
// @Summary 获取用户钱包信息
// @Description 获取当前登录用户的钱包信息
// @Tags 用户/钱包
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /v1/user/wallet [get]
func GetWallet(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var wallet walletModel.UserWallet
	if err := global.APP_DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		// 钱包不存在则创建
		newWallet := walletModel.UserWallet{
			UserID:        userID.(uint),
			Balance:       0,
			Frozen:        0,
			TotalRecharge: 0,
			TotalExpense:  0,
		}
		if err := global.APP_DB.Create(&newWallet).Error; err != nil {
			global.APP_LOG.Error("创建钱包失败", zap.Error(err))
			c.JSON(500, gin.H{"code": 500, "message": "获取钱包信息失败"})
			return
		}
		wallet = newWallet
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    wallet,
	})
}

// GetWalletTransactions 获取钱包交易记录
// @Summary 获取钱包交易记录
// @Description 获取当前登录用户的交易记录
// @Tags 用户/钱包
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param type query string false "交易类型"
// @Success 200 {object} common.Response
// @Router /v1/user/wallet/transactions [get]
func GetWalletTransactions(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未授权"})
		return
	}

	// 分页参数
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := parseUint(p); err == nil && parsed > 0 {
			page = int(parsed)
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if parsed, err := parseUint(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = int(parsed)
		}
	}

	var transactions []walletModel.WalletTransaction
	var total int64

	query := global.APP_DB.Model(&walletModel.WalletTransaction{}).Where("user_id = ?", userID)

	// 交易类型筛选
	if txType := c.Query("type"); txType != "" {
		query = query.Where("type = ?", txType)
	}

	// 查询总数
	query.Count(&total)

	// 查询列表
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		global.APP_LOG.Error("获取交易记录失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": "获取交易记录失败"})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":     transactions,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// parseUint 解析uint
func parseUint(s string) (uint, error) {
	var result uint
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
