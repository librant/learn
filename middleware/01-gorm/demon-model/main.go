package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	log.SetFlags(log.Lshortfile)

	db, err := gorm.Open("mysql",
		"root:123456@(127.0.0.1:3306)/db6?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	db.AutoMigrate(&User{})
}
