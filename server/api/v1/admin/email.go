package admin

import (
	"net/http"
	"oneclickvirt/model/common"
	"oneclickvirt/service/email"

	"github.com/gin-gonic/gin"
)

// EmailApi 邮件API
type EmailApi struct{}

// SendTestEmail 发送测试邮件
// @Summary 发送测试邮件
// @Description 发送测试邮件
// @Tags 邮件管理
// @Accept json
// @Produce json
// @Param request body object{to=string} true "发送测试邮件请求"
// @Success 200 {object} common.Response
// @Router /admin/email/test [post]
func (api *EmailApi) SendTestEmail(c *gin.Context) {
	var req struct {
		To string `json:"to" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Code: 400,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	service := email.NewEmailService()
	if err := service.SendTestEmail(req.To); err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Code: 500,
			Msg:  "发送测试邮件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Msg:  "发送成功",
	})
}

// GetEmailConfig 获取邮件配置
// @Summary 获取邮件配置
// @Description 获取邮件配置
// @Tags 邮件管理
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=object}
// @Router /admin/email/config [get]
func (api *EmailApi) GetEmailConfig(c *gin.Context) {
	service := email.NewEmailService()

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Data: gin.H{
			"smtp_host":  service.SMTPHost,
			"smtp_port":  service.SMTPPort,
			"from":      service.From,
			"from_name": service.FromName,
			"use_ssl":   service.UseSSL,
			"password":  "***", // 隐藏密码
		},
		Msg: "获取成功",
	})
}

// UpdateEmailConfig 更新邮件配置
// @Summary 更新邮件配置
// @Description 更新邮件配置
// @Tags 邮件管理
// @Accept json
// @Produce json
// @Param request body object{smtp_host=string,smtp_port=string,from=string,from_name=string,password=string,use_ssl=bool} true "更新邮件配置请求"
// @Success 200 {object} common.Response
// @Router /admin/email/config [put]
func (api *EmailApi) UpdateEmailConfig(c *gin.Context) {
	var req struct {
		SMTPHost  string `json:"smtp_host" binding:"required"`
		SMTPPort  string `json:"smtp_port" binding:"required"`
		From      string `json:"from" binding:"required"`
		FromName  string `json:"from_name" binding:"required"`
		Password  string `json:"password" binding:"required"`
		UseSSL    bool   `json:"use_ssl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Code: 400,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	// TODO: 保存配置到数据库或配置文件

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Msg:  "更新成功",
	})
}
