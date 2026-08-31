package models

// User represents an application user with credentials and role information.
type User struct {
	Id       int    `gorm:"primaryKey;column:id"`
	Username string `gorm:"column:username;unique" json:"username" binding:"required"`
	Email    string `gorm:"column:email;unique" json:"email" binding:"required,email"`
	Password string `gorm:"column:password" json:"password" binding:"required,min=6"`
	Role     string `gorm:"column:role" json:"role" binding:"required"`
}
