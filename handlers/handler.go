// Package handlers provides HTTP request handlers for the application.
package handlers

import "gorm.io/gorm"

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	Db *gorm.DB
}

// SetHandler initializes and returns a new Handler instance.
func SetHandler(db *gorm.DB) *Handler {
	return &Handler{Db: db}
}
