package user

import (
	"strconv"

	"oneclickvirt/middleware"
	"oneclickvirt/model/common"
	paymentSvc "oneclickvirt/service/payment"
	productSvc "oneclickvirt/service/product"

	"github.com/gin-gonic/gin"
)

// ListProducts 产品列表（用户端，已上架）
// @Router /user/products [get]
func ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))
	instanceType := c.Query("type")
	svc := productSvc.NewService()
	list, total, err := svc.ListProducts(page, pageSize, instanceType)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccessWithPagination(c, list, total, page, pageSize)
}

// GetProduct 产品详情（用户端）
// @Router /user/products/:id [get]
func GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}
	svc := productSvc.NewService()
	p, err := svc.Get(uint(id))
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccess(c, p)
}

// BuyProductRequest 购买请求体
type BuyProductRequest struct {
	Method     string `json:"method" binding:"required,oneof=balance epay"`
	NodeID     uint   `json:"nodeId"`
	ImageID    uint   `json:"imageId"`
	PeriodDays int    `json:"periodDays"`
	RenewalOf  *uint  `json:"renewalOf"`
	BaseURL    string `json:"baseURL"`
}

// BuyProduct 购买产品（余额或易支付）
// @Router /user/products/:id/buy [post]
func BuyProduct(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "无效的产品ID"))
		return
	}
	var req BuyProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeValidationError, "参数错误: "+err.Error()))
		return
	}
	svc := productSvc.NewService()
	res, err := svc.Buy(userID, uint(id), productSvc.BuyRequest{
		Method:        req.Method,
		NodeID:        req.NodeID,
		ImageID:       req.ImageID,
		PeriodDays:    req.PeriodDays,
		RenewalOf:     req.RenewalOf,
		ClientBaseURL: req.BaseURL,
	})
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	data := gin.H{
		"tradeNo": res.Order.TradeNo,
		"orderId": res.Order.ID,
		"status":  res.Order.Status,
	}
	if res.EpayParams != nil {
		data["epay"] = res.EpayParams
		data["epayUrl"] = paymentSvc.EpaySubmitURL()
	}
	common.ResponseSuccess(c, data, "下单成功")
}

// ListMyOrders 我的订单
// @Router /user/orders [get]
func ListMyOrders(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		common.ResponseWithError(c, common.NewError(common.CodeUnauthorized, err.Error()))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	svc := productSvc.NewService()
	list, total, err := svc.MyOrders(userID, page, pageSize)
	if err != nil {
		common.ResponseWithError(c, common.ClassifyError(err))
		return
	}
	common.ResponseSuccessWithPagination(c, list, total, page, pageSize)
}
