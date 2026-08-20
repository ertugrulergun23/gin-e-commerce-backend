package models

import "time"

type Order struct {
	Id         int       `gorm:"column:id;primaryKey"`
	Owner_id   int       `gorm:"column:owner_id"`
	Status     string    `gorm:"column:status" json:"status"`
	created_at time.Time `gorm:"column:created_at"`
}

type Order_Item struct {
	Id         int `gorm:"column:id;primaryKey;primaryKey"`
	Order_id   int `gorm:"column:order_id"`
	Product_id int `gorm:"column:product_id"`
	Quantity   int `gorm:"column:quantity"`
}
