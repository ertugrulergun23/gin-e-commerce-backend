// Package routers defines and registers HTTP routes for the application.
package routers

import (
	"ecommerce/handlers"
	"ecommerce/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouters configures the HTTP router, middleware, and API endpoints.
func SetupRouters(db *gorm.DB) *gin.Engine {
	// Initialize default Gin router instance
	router := gin.Default()

	// Initialize handler with database connection
	h := handlers.SetHandler(db)

	// Authentication endpoint with JWT
	router.POST("/login", h.Login)

	// Protected API route group requiring authentication middleware
	api := router.Group("/api")
	api.Use(middleware.AuthRequired)
	{
		// User model endpoints
		api.GET("/user", h.GetUser)
		api.PUT("/user/update", h.UpdateUser)
		api.DELETE("/user/delete", h.DeleteUser)

		// Product model endpoints
		api.POST("/product/create", h.CreateProduct)
		api.PUT("/product/:id/update", h.UpdateProduct)
		api.DELETE("/product/:id/delete", h.DeleteProduct)

		// Comment model endpoints
		api.POST("/comment/create/:product_id", h.CreateComment)
		api.PUT("/comment/:id/update", h.UpdateComment)
		api.DELETE("/comment/:id/delete", h.DeleteComment)

		// Cart model endpoints
		api.GET("/cart", h.GetUserCart)
		api.POST("/cart/add", h.AddProductToCart)
		api.PUT("/cart/:id/update", h.UpdateCart)
		api.DELETE("/cart/:id/delete", h.DeleteCart)

		// Order model endpoints
		api.GET("/order", h.GetUserOrders)
		api.POST("/order/cart", h.CreateOrderFromCart)
		api.POST("/order/create", h.CreateOrder)

		// Admin-only order endpoint
		api.PUT("/order/:id/update", h.UpdateOrderStatus)
	}

	// Public user registration endpoint
	router.POST("/user/create", h.CreateUser)

	// Public get all products with filters endpoint
	router.GET("/products", h.GetProducts)

	// Public get product with filters endpoint
	router.GET("/product/:id", h.GetProduct)

	// Public get comment endpoint
	router.GET("/comment/:id", h.GetComment)

	// Public get comments with filters endpoint
	router.GET("/comments", h.GetComments)

	// Test health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
