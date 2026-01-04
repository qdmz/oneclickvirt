package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SiteConfig 站点配置模型
type SiteConfig struct {
	ID          uint   `gorm:"primaryKey"`
	Key         string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Value       string `gorm:"type:text"`
	Type        string `gorm:"type:varchar(20);not null"`
	Description string `gorm:"type:varchar(255)"`
}

func main() {
	// 确保当前目录有数据文件
	if _, err := os.Stat("data.db"); os.IsNotExist(err) {
		log.Fatal("data.db not found!")
	}

	// 使用CGO_ENABLED=1编译
	os.Setenv("CGO_ENABLED", "1")

	// 连接数据库
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&SiteConfig{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// 检查contact_phone配置项是否存在
	var config SiteConfig
	result := db.Where("key = ?", "contact_phone").First(&config)

	if result.Error != nil {
		// 如果不存在，创建它
		config = SiteConfig{
			Key:         "contact_phone",
			Value:       "888-888-8888",
			Type:        "string",
			Description: "联系电话",
		}
		if err := db.Create(&config).Error; err != nil {
			log.Fatal("Failed to create contact_phone config: ", err)
		}
		fmt.Println("contact_phone config created successfully!")
	} else {
		// 如果存在，更新它的值
		config.Value = "888-888-8888"
		if err := db.Save(&config).Error; err != nil {
			log.Fatal("Failed to update contact_phone config: ", err)
		}
		fmt.Println("contact_phone config updated successfully!")
	}

	// 再次查询确认
	if err := db.Where("key = ?", "contact_phone").First(&config).Error; err != nil {
		log.Fatal("Failed to verify contact_phone config: ", err)
	}
	fmt.Printf("Final contact_phone config: Key = %s, Value = %s\n", config.Key, config.Value)
}
