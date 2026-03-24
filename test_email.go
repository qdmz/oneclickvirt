package main

import (
	"fmt"
	"log"

	"oneclickvirt/config"
	"oneclickvirt/global"
	"oneclickvirt/service/auth"

	"github.com/spf13/viper"
)

func main() {
	// 加载配置
	viper.SetConfigFile("server/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 初始化全局配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 创建AuthService实例
	authService := &auth.AuthService{}

	// 测试邮件发送
	to := "test@example.com"
	subject := "测试邮件"
	body := "这是一封测试邮件，用于验证邮件发送功能是否正常工作。"

	fmt.Printf("正在发送测试邮件到 %s...\n", to)
	err := authService.SendVerifyCode("email", to)
	if err != nil {
		fmt.Printf("发送邮件失败: %v\n", err)
	} else {
		fmt.Println("邮件发送成功！")
	}
}
