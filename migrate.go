package main

import (
	"fmt"
	"log"

	"oneclickvirt/global"
	"oneclickvirt/initialize"
	"oneclickvirt/model/config"
)

func main() {
	// 加载配置
	global.APP_VP = initialize.Viper()
	
	// 初始化数据库连接
	dbManager := initialize.GetDatabaseManager()
	mysqlConfig := config.MysqlConfig{
		Path:         global.APP_CONFIG.Mysql.Path,
		Port:         global.APP_CONFIG.Mysql.Port,
		Config:       global.APP_CONFIG.Mysql.Config,
		Dbname:       global.APP_CONFIG.Mysql.Dbname,
		Username:     global.APP_CONFIG.Mysql.Username,
		Password:     global.APP_CONFIG.Mysql.Password,
		MaxIdleConns: global.APP_CONFIG.Mysql.MaxIdleConns,
		MaxOpenConns: global.APP_CONFIG.Mysql.MaxOpenConns,
		LogMode:      global.APP_CONFIG.Mysql.LogMode,
		LogZap:       global.APP_CONFIG.Mysql.LogZap,
		MaxLifetime:  global.APP_CONFIG.Mysql.MaxLifetime,
		AutoCreate:   global.APP_CONFIG.Mysql.AutoCreate,
	}

	db, err := dbManager.Initialize(mysqlConfig)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	fmt.Println("数据库连接成功")

	// 执行数据库迁移
	fmt.Println("开始数据库表结构自动迁移")
	initialize.RegisterTables(db)
	fmt.Println("数据库表结构迁移完成")

	fmt.Println("数据库迁移成功")
}
