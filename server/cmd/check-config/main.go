package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"oneclickvirt/model/site"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	// 获取所有配置项
	var configs []site.SiteConfig
	if err := db.Find(&configs).Error; err != nil {
		log.Fatal("Failed to get configs: ", err)
	}

	fmt.Println("All site configs:")
	for _, config := range configs {
		fmt.Printf("Key: %s, Value: %s\n", config.Key, config.Value)
	}

	// 检查contact_phone配置项
	var contactPhoneConfig site.SiteConfig
	if err := db.Where("key = ?", "contact_phone").First(&contactPhoneConfig).Error; err != nil {
		fmt.Println("\ncontact_phone config not found!")
		
		// 如果不存在，创建它
		fmt.Println("Creating contact_phone config...")
		contactPhoneConfig = site.SiteConfig{
			Key:         "contact_phone",
			Value:       "888-888-8888",
			Type:        "string",
			Description: "联系电话",
		}
		if err := db.Create(&contactPhoneConfig).Error; err != nil {
			log.Fatal("Failed to create contact_phone config: ", err)
		}
		fmt.Println("contact_phone config created successfully!")
	} else {
		fmt.Printf("\ncontact_phone config found: Value = %s\n", contactPhoneConfig.Value)
	}
}
