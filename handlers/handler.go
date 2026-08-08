package handlers

import "gorm.io/gorm"

type Handler struct {
	Db *gorm.DB
}

func SetHandler(db *gorm.DB) *Handler {
	return &Handler{Db: db}
}
