package models

type User struct {
	Id       int    `gorm:"primarykey;column:id" json:"id"`
	Username string `gorm:"column:username" json:"username" binding:"required"`
	Email    string `gorm:"column:email" json:"email" binding:"required,email"`
	Password string `gorm:"column:password" json:"password" binding:"required,min=6"`
	Seller   bool   `gorm:"column:seller" json:"seller" binding:"required"`
}
