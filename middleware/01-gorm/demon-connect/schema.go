package main

// UserInfo 用户信息
type UserInfo struct {
	ID   uint   `gorm:"column:f_id"`
	Name string `gorm:"column:f_name"`
}
