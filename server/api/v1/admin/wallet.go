package admin

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	walletSvc "oneclickvirt/service/wallet"

	"github.com/gin-gonic/gin"
)

// AdjustBalanceRequest 管理员调整用户余额请求
type AdjustBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required"` // 元，可为负数表示扣减
	Remark string  `json:"remark"`
}

// AdjustUserBalance 管理员调整用户余额
// @Router /admin/users/:id/balance [post]
func AdjustUserBalance(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		common.ResponseWithError(c, common.NewError(common.CodeForbidden, "需要超管权限"))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的用户ID"))
		return
	}
	var req AdjustBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}
	cents := int64(req.Amount*100 + 0.5)
	if req.Amount < 0 {
		cents = int64(req.Amount*100 - 0.5)
	}
	svc := walletSvc.NewService()
	balance, err := svc.AdjustBalance(uint(id), cents, req.Remark)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, gin.H{"balance": balance}, "余额已调整")
}

// ListRechargeOrdersAdmin 管理端充值/余额变动流水
// @Router /admin/recharge-orders [get]
func ListRechargeOrdersAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	userID, _ := strconv.ParseUint(c.DefaultQuery("userId", "0"), 10, 64)
	svc := walletSvc.NewService()
	list, total, err := svc.AdminListRechargeOrders(uint(userID), page, pageSize)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccessWithPagination(c, list, total, page, pageSize)
}

// ListOrdersAdmin 管理端商品订单列表
// @Router /admin/orders [get]
func ListOrdersAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	userID, _ := strconv.ParseUint(c.DefaultQuery("userId", "0"), 10, 64)
	status := c.Query("status")
	svc := walletSvc.NewService()
	list, total, err := svc.AdminListOrders(uint(userID), status, page, pageSize)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccessWithPagination(c, list, total, page, pageSize)
}
