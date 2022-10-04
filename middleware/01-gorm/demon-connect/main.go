package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// 1 配置 mysql 链接数据库
	userName := "root"
	password := "123456"
	host := "127.0.0.1"
	port := 3306
	dbName := "test"
	timeout := "10s"

	// 2 拼接 dsn 参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		userName, password, host, port, dbName, timeout)

	// 3 链接 mysql, 获得 DB 类型实例
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	// 创建表 自动迁移
	db.AutoMigrate(&UserInfo{})
}
