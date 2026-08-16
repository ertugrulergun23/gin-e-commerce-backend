package routers

import (
	"ecommerce/handlers"
	"ecommerce/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouters(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	h := handlers.SetHandler(db)

	// authentication with JWT
	router.POST("/login", h.Login)

	api := router.Group("/api")
	api.Use(middleware.AuthRequired)
	{
		// user model endpoints
		api.GET("/user", h.GetUser)
		api.PUT("/user/update", h.UpdateUser)
		api.DELETE("user/delete", h.DeleteUser)
		// product model enpoints
		api.POST("/product/create", h.CreateProduct)
		api.PUT("/product/:id/update", h.UpdateProduct)
		api.DELETE("/product/:id/delete", h.DeleteProduct)
		// comment model enpoints
		api.POST("/comment/create/:product_id", h.CreateComment)
		api.PUT("/comment/:id/update", h.UpdateComment)
		api.DELETE("/comment/:id/delete", h.DeleteComment)
		// cart model enpoint
		api.GET("/cart", h.GetUserCart)
		api.POST("/cart/add", h.AddProductToCart)
		api.PUT("/cart/:id/update", h.UpdateCart)
		api.DELETE("/cart/:id/delete", h.DeleteCart)

	}

	// user create endpoint
	router.POST("/user/create", h.CreateUser)

	// get all products with filters
	router.GET("/products", h.GetProducts)

	// Public get product with filters endpoint
	router.GET("/product/:id", h.GetProduct)

	// Public get comment endpoint
	router.GET("/comment/:id", h.GetComment)

	// Public get comments with filters endpoint
	router.GET("/comments", h.GetComments)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
