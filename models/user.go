package models

type User struct {
	Id       int    `gorm:"primaryKey;column:id"`
	Username string `gorm:"column:username;unique" json:"username" binding:"required"`
	Email    string `gorm:"column:email;unique" json:"email" binding:"required,email"`
	Password string `gorm:"column:password" json:"password" binding:"required,min=6"`
	Seller   *bool  `gorm:"column:seller" json:"seller" binding:"required"`
}
