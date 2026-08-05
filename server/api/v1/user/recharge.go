package user

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	paymentSvc "oneclickvirt/service/payment"
	walletSvc "oneclickvirt/service/wallet"

	"github.com/gin-gonic/gin"
)

// GetWallet 获取钱包余额
// @Router /user/wallet [get]
func GetWallet(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}
	svc := walletSvc.NewService()
	balance, err := svc.GetBalance(userID)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, gin.H{"balance": balance})
}

// CreateRechargeRequest 发起充值请求
type CreateRechargeRequest struct {
	Amount  float64 `json:"amount" binding:"required,gt=0"` // 元
	BaseURL string  `json:"baseURL"`                        // 前端传入，用于拼接回调地址
}

// CreateRecharge 创建余额充值订单（返回易支付跳转参数）
// @Router /user/recharge [post]
func CreateRecharge(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}
	var req CreateRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}
	cents := int64(req.Amount * 100 + 0.5)
	svc := walletSvc.NewService()
	order, params, err := svc.CreateRechargeOrder(userID, cents, req.BaseURL)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, gin.H{
		"tradeNo": order.TradeNo,
		"amount":  order.Amount,
		"epay":    params,
		"epayUrl": paymentSvc.EpaySubmitURL(),
	}, "充值订单已创建，请前往支付")
}

// ListRechargeOrders 我的充值订单
// @Router /user/recharge/orders [get]
func ListRechargeOrders(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	svc := walletSvc.NewService()
	list, total, err := svc.ListUserRechargeOrders(userID, page, pageSize)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccessWithPagination(c, list, total, page, pageSize)
}
