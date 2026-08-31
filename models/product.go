package models

// Product represents an item listed for sale by a seller.
type Product struct {
	Id        int     `gorm:"primaryKey;column:id" json:"id"`
	Name      string  `gorm:"column:name" json:"name" binding:"required"`
	Price     float64 `gorm:"column:price" json:"price" binding:"required,gt=0"`
	Stock     int     `gorm:"column:stock" json:"stock" binding:"gte=0"`
	Point     float64 `gorm:"column:point" json:"point"`
	Seller_id int     `gorm:"column:seller_id" json:"seller_id" binding:"required"`
}
