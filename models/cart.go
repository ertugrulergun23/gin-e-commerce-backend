package models

import "time"

type Cart struct {
	Id         int       `gorm:"primarykey;column:id" json:"id"`
	User_id    int       `gorm:"column:user_id" json:"user_id"`
	Product_id int       `gorm:"column:product_id" json:"product_id"`
	Quantity   int       `gorm:"column:quantity;default:1" json:"quantity"`
	Created_at time.Time `gorm:"column:created_at"`
}
