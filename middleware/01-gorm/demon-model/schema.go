package main

import (
	"time"

	"database/sql"
	"github.com/jinzhu/gorm"
)

// UserInfo 用户信息
type UserInfo struct {
	ID     uint
	Name   string
	Gender string
	Hobby  string
}

// User 用户模型
type User struct {
	gorm.Model
	Name      string
	Age       sql.NullString
	Birthday  *time.Time
	Email     string  `gorm:"type:varchar(100);unique_index"`
	Role      string  `gorm:"size:255"`
	MemberNum *string `gorm:"unique;not null"`
	Num       int     `gorm:"AUTO_INCREMENT"`
	Addr      string  `gorm:"index:addr)"`
	IgnoreMe  int     `gorm:"-"`
}
