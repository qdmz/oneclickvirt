package main

import (
	"fmt"
	"oneclickvirt/core"
	"oneclickvirt/global"
	"oneclickvirt/initialize"
	"oneclickvirt/model/order"
	"oneclickvirt/model/product"
	"oneclickvirt/model/redemption"
	"oneclickvirt/source"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// 初始化核心组件
	global.APP_VP = core.Viper()
	global.APP_LOG = core.Zap()
	zap.ReplaceGlobals(global.APP_LOG)

	fmt.Println("正在连接数据库...")
	// 连接数据库
	global.APP_DB = initialize.Gorm()
	if global.APP_DB == nil {
		fmt.Println("数据库连接失败")
		return
	}

	fmt.Println("数据库连接成功")

	// 按照外键约束顺序删除数据
	// 1. 先删除支付记录（依赖订单）
	fmt.Println("正在删除所有支付记录...")
	paymentResult := global.APP_DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&order.PaymentRecord{})
	if paymentResult.Error != nil {
		fmt.Printf("删除支付记录失败: %v\n", paymentResult.Error)
		return
	}
	fmt.Printf("成功删除 %d 个支付记录\n", paymentResult.RowsAffected)

	// 2. 然后删除订单（依赖产品）
	fmt.Println("正在删除所有订单记录...")
	orderResult := global.APP_DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&order.Order{})
	if orderResult.Error != nil {
		fmt.Printf("删除订单失败: %v\n", orderResult.Error)
		return
	}
	fmt.Printf("成功删除 %d 个订单\n", orderResult.RowsAffected)

	// 3. 先删除兑换码使用记录（依赖兑换码）
	fmt.Println("正在删除所有兑换码使用记录...")
	redemptionUsageResult := global.APP_DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&redemption.RedemptionCodeUsage{})
	if redemptionUsageResult.Error != nil {
		fmt.Printf("删除兑换码使用记录失败: %v\n", redemptionUsageResult.Error)
		return
	}
	fmt.Printf("成功删除 %d 个兑换码使用记录\n", redemptionUsageResult.RowsAffected)

	// 4. 然后删除兑换码（依赖产品）
	fmt.Println("正在删除所有兑换码...")
	redemptionResult := global.APP_DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&redemption.RedemptionCode{})
	if redemptionResult.Error != nil {
		fmt.Printf("删除兑换码失败: %v\n", redemptionResult.Error)
		return
	}
	fmt.Printf("成功删除 %d 个兑换码\n", redemptionResult.RowsAffected)

	// 5. 最后删除产品
	fmt.Println("正在删除所有产品...")
	productResult := global.APP_DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&product.Product{})
	if productResult.Error != nil {
		fmt.Printf("删除产品失败: %v\n", productResult.Error)
		return
	}
	fmt.Printf("成功删除 %d 个产品\n", productResult.RowsAffected)

	// 重新生成产品
	fmt.Println("正在重新生成产品...")
	source.InitSeedData()
	fmt.Println("产品重新生成完成")

	// 验证生成的产品数量
	var count int64
	global.APP_DB.Model(&product.Product{}).Count(&count)
	fmt.Printf("当前共有 %d 个产品\n", count)

	fmt.Println("操作完成")
}
