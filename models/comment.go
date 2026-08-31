package models

import "time"

// Comment represents a review and rating left on a product.
type Comment struct {
	Id         int       `gorm:"primaryKey; column:id" json:"id"`
	Owner_id   int       `gorm:"column:owner_id" json:"owner_id" binding:"required"`
	Product_id int       `gorm:"column:product_id" json:"product_id" binding:"required"`
	Content    string    `gorm:"column:content" json:"content" binding:"required"`
	Point      int       `gorm:"column:point" json:"point" binding:"required"`
	Created_at time.Time `gorm:"column:created_at" json:"created_at"`
}
