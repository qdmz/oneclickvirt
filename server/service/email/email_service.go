package email

import (
	"fmt"
	"net/smtp"
	"oneclickvirt/global"
	"strings"

	"go.uber.org/zap"
)

// EmailService 邮件服务
type EmailService struct {
	// SMTP服务器地址
	SMTPHost string
	// SMTP端口
	SMTPPort string
	// 发件人邮箱
	From string
	// 发件人名称
	FromName string
	// 授权密码
	Password string
	// 是否使用SSL
	UseSSL bool
}

// NewEmailService 创建邮件服务
func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost: "smtp.qq.com",
		SMTPPort: "587", // 使用TLS端口
		From:     "qdmz@vip.qq.com",
		FromName: "OneClickVirt",
		Password: "jehpkshzykrlbjjf",
		UseSSL:   false, // 使用TLS
	}
}

// SendEmail 发送邮件
func (s *EmailService) SendEmail(to []string, subject, body string) error {
	// 构建邮件内容
	msg := s.buildMessage(to, subject, body)

	// 发送邮件
	var err error
	if s.UseSSL {
		err = s.sendWithSSL(to, msg)
	} else {
		err = s.sendWithTLS(to, msg)
	}

	if err != nil {
		global.APP_LOG.Error("发送邮件失败",
			zap.Strings("to", to),
			zap.String("subject", subject),
			zap.Error(err))
		return err
	}

	global.APP_LOG.Info("发送邮件成功",
		zap.Strings("to", to),
		zap.String("subject", subject))

	return nil
}

// buildMessage 构建邮件内容
func (s *EmailService) buildMessage(to []string, subject, body string) string {
	var msg strings.Builder

	// 发件人
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.FromName, s.From))

	// 收件人
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))

	// 主题
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))

	// MIME版本
	msg.WriteString("MIME-Version: 1.0\r\n")

	// 内容类型
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")

	// 空行
	msg.WriteString("\r\n")

	// 邮件正文
	msg.WriteString(body)

	return msg.String()
}

// sendWithSSL 使用SSL发送邮件
func (s *EmailService) sendWithSSL(to []string, msg string) error {
	// 构建SMTP地址
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)

	// 认证信息
	auth := smtp.PlainAuth("", s.From, s.Password, s.SMTPHost)

	// 发送邮件
	return smtp.SendMail(addr, auth, s.From, to, []byte(msg))
}

// sendWithTLS 使用TLS发送邮件
func (s *EmailService) sendWithTLS(to []string, msg string) error {
	// 构建SMTP地址
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)

	// 认证信息
	auth := smtp.PlainAuth("", s.From, s.Password, s.SMTPHost)

	// 发送邮件
	return smtp.SendMail(addr, auth, s.From, to, []byte(msg))
}

// SendTestEmail 发送测试邮件
func (s *EmailService) SendTestEmail(to string) error {
	subject := "OneClickVirt 测试邮件"
	body := `这是一封测试邮件。

如果您收到这封邮件，说明邮件服务配置正确。

---
OneClickVirt
https://github.com/qdmz/oneclickvirt`

	return s.SendEmail([]string{to}, subject, body)
}

// SendOrderEmail 发送订单邮件
func (s *EmailService) SendOrderEmail(to, orderNo string, amount float64) error {
	subject := fmt.Sprintf("订单支付成功 - %s", orderNo)
	body := fmt.Sprintf(`尊敬的用户：

您的订单 %s 已支付成功！

订单金额：%.2f 元

感谢您的使用！

---
OneClickVirt
https://github.com/qdmz/oneclickvirt`, orderNo, amount)

	return s.SendEmail([]string{to}, subject, body)
}

// SendWelcomeEmail 发送欢迎邮件
func (s *EmailService) SendWelcomeEmail(to, username string) error {
	subject := "欢迎注册 OneClickVirt"
	body := fmt.Sprintf(`尊敬的 %s：

欢迎注册 OneClickVirt！

OneClickVirt 是一个一键虚拟化管理平台，为您提供便捷的虚拟机管理服务。

如果您有任何问题，请随时联系我们。

---
OneClickVirt
https://github.com/qdmz/oneclickvirt`, username)

	return s.SendEmail([]string{to}, subject, body)
}

// SendResetPasswordEmail 发送密码重置邮件
func (s *EmailService) SendResetPasswordEmail(to, username, resetURL string) error {
	subject := "密码重置请求"
	body := fmt.Sprintf(`尊敬的 %s：

我们收到了您的密码重置请求。

请点击以下链接重置您的密码：
%s

如果您没有请求密码重置，请忽略此邮件。

---
OneClickVirt
https://github.com/qdmz/oneclickvirt`, username, resetURL)

	return s.SendEmail([]string{to}, subject, body)
}

// SendActivationEmail 发送激活邮件
func (s *EmailService) SendActivationEmail(to, token, frontendURL string) error {
	subject := "账户激活"
	activationURL := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)
	body := fmt.Sprintf(`尊敬的用户：

感谢您注册 OneClickVirt！

请点击以下链接激活您的账户：
%s

如果您没有注册，请忽略此邮件。

---
OneClickVirt
https://github.com/qdmz/oneclickvirt`, activationURL)

	return s.SendEmail([]string{to}, subject, body)
}
