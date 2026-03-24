package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net"

	"oneclickvirt/global"
	"go.uber.org/zap"
)

// EmailService 邮件服务
type EmailService struct{}

// NewEmailService 创建邮件服务实例
func NewEmailService() *EmailService {
	return &EmailService{}
}

// sendEmail 发送邮件通用方法
func (s *EmailService) sendEmail(to, subject, htmlBody string) error {
	config := global.APP_CONFIG.Auth
	if config.EmailSMTPHost == "" {
		return fmt.Errorf("邮件服务未配置")
	}

	if config.EmailUsername == "" || config.EmailPassword == "" {
		return fmt.Errorf("邮件服务用户名或密码未配置")
	}

	auth := smtp.PlainAuth("", config.EmailUsername, config.EmailPassword, config.EmailSMTPHost)
	msg := "To: " + to + "\r\n" +
		"From: " + config.EmailUsername + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"MIME-Version: 1.0\r\n" +
		"\r\n" +
		htmlBody

	// 处理安全连接
	addr := fmt.Sprintf("%s:%d", config.EmailSMTPHost, config.EmailSMTPPort)
	
	// 建立连接
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		global.APP_LOG.Error("连接邮件服务器失败", zap.Error(err), zap.String("addr", addr))
		return fmt.Errorf("连接邮件服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, config.EmailSMTPHost)
	if err != nil {
		global.APP_LOG.Error("创建SMTP客户端失败", zap.Error(err), zap.String("host", config.EmailSMTPHost))
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 启用TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.EmailSMTPHost,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		global.APP_LOG.Error("启用TLS失败", zap.Error(err))
		return fmt.Errorf("启用TLS失败: %v", err)
	}

	// 认证
	if err := client.Auth(auth); err != nil {
		global.APP_LOG.Error("邮件服务器认证失败", zap.Error(err), zap.String("username", config.EmailUsername))
		return fmt.Errorf("邮件服务器认证失败: %v", err)
	}

	// 设置发件人
	if err := client.Mail(config.EmailUsername); err != nil {
		global.APP_LOG.Error("设置发件人失败", zap.Error(err), zap.String("from", config.EmailUsername))
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	if err := client.Rcpt(to); err != nil {
		global.APP_LOG.Error("设置收件人失败", zap.Error(err), zap.String("to", to))
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		global.APP_LOG.Error("准备发送数据失败", zap.Error(err))
		return fmt.Errorf("准备发送数据失败: %v", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		global.APP_LOG.Error("写入邮件数据失败", zap.Error(err))
		return fmt.Errorf("写入邮件数据失败: %v", err)
	}

	if err := w.Close(); err != nil {
		global.APP_LOG.Error("关闭写入流失败", zap.Error(err))
		return fmt.Errorf("关闭写入流失败: %v", err)
	}

	if err := client.Quit(); err != nil {
		global.APP_LOG.Error("关闭SMTP连接失败", zap.Error(err))
		return fmt.Errorf("关闭SMTP连接失败: %v", err)
	}

	global.APP_LOG.Info("邮件发送成功", zap.String("to", to), zap.String("subject", subject))
	return nil
}

// SendVerificationCode 发送注册验证码
func (s *EmailService) SendVerificationCode(email, code string) error {
	subject := "注册验证码"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <div style="background: #f8f9fa; border-radius: 8px; padding: 30px;">
    <h2 style="color: #16a34a;">注册验证码</h2>
    <p>您好，您的注册验证码是：</p>
    <div style="font-size: 32px; font-weight: bold; color: #16a34a; padding: 20px; text-align: center; background: #e8f5e9; border-radius: 8px; letter-spacing: 4px;">%s</div>
    <p style="color: #666;">验证码5分钟内有效，请勿泄露给他人。</p>
  </div>
</body>
</html>`, code)
	return s.sendEmail(email, subject, body)
}

// SendActivationEmail 发送激活链接
func (s *EmailService) SendActivationEmail(email, token, frontendURL string) error {
	subject := "邮箱激活"
	activateURL := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <div style="background: #f8f9fa; border-radius: 8px; padding: 30px;">
    <h2 style="color: #16a34a;">邮箱激活</h2>
    <p>您好，感谢您的注册！请点击以下链接激活您的邮箱：</p>
    <div style="text-align: center; padding: 20px;">
      <a href="%s" style="display: inline-block; padding: 12px 30px; background: #16a34a; color: #fff; text-decoration: none; border-radius: 6px; font-size: 16px;">激活邮箱</a>
    </div>
    <p style="color: #666;">或者复制以下链接到浏览器：</p>
    <p style="color: #409eff; word-break: break-all;">%s</p>
    <p style="color: #999; font-size: 12px;">此链接有效期为24小时。</p>
  </div>
</body>
</html>`, activateURL, activateURL)
	return s.sendEmail(email, subject, body)
}

// SendPasswordResetEmail 发送密码重置链接
func (s *EmailService) SendPasswordResetEmail(email, token, frontendURL string) error {
	subject := "密码重置"
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <div style="background: #f8f9fa; border-radius: 8px; padding: 30px;">
    <h2 style="color: #e6a23c;">密码重置</h2>
    <p>您好，您正在重置密码，请点击以下链接：</p>
    <div style="text-align: center; padding: 20px;">
      <a href="%s" style="display: inline-block; padding: 12px 30px; background: #e6a23c; color: #fff; text-decoration: none; border-radius: 6px; font-size: 16px;">重置密码</a>
    </div>
    <p style="color: #666;">或者复制以下链接到浏览器：</p>
    <p style="color: #409eff; word-break: break-all;">%s</p>
    <p style="color: #999; font-size: 12px;">此链接有效期为24小时。如果这不是您本人的操作，请忽略此邮件。</p>
  </div>
</body>
</html>`, resetURL, resetURL)
	return s.sendEmail(email, subject, body)
}

// SendWelcomeEmail 发送欢迎邮件
func (s *EmailService) SendWelcomeEmail(email, username string) error {
	subject := "欢迎注册"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <div style="background: #f8f9fa; border-radius: 8px; padding: 30px;">
    <h2 style="color: #16a34a;">欢迎注册！🎉</h2>
    <p>您好 <strong>%s</strong>，欢迎加入！</p>
    <p>您的账号已成功创建并激活。现在您可以：</p>
    <ul>
      <li>登录您的账户</li>
      <li>管理您的实例</li>
      <li>查看您的资源配额</li>
    </ul>
    <p>如有任何问题，请随时联系管理员。</p>
  </div>
</body>
</html>`, username)
	return s.sendEmail(email, subject, body)
}