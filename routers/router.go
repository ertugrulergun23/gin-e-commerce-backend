package routers

import (
	"ecommerce/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouters(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	h := handlers.SetHandler(db)

	// users database CRUD endpoints
	router.POST("/user/create", h.CreateUser)
	router.GET("user/:id", h.GetUser)
	router.PUT("/user/:id/update", h.UpdateUser)
	router.DELETE("/user/:id/delete", h.DeleteUser)

	// products database CRUD endpoints
	router.GET("/products", h.GetProducts)
	router.GET("/product/:id", h.GetProduct)
	router.POST("/product/create", h.CreateProduct)
	router.PUT("/product/:id/update", h.UpdateProduct)
	router.DELETE("/product/:id/delete", h.DeleteProduct)

	// comments database CRUD enpoints
	router.GET("/comment/:id", h.GetComment)
	router.POST("/comment/create", h.CreateComment)
	router.PUT("/comment/:id/update", h.UpdateComment)
	router.DELETE("/comment/:İd/delete", h.DeleteComment)
	router.GET("/comments", h.GetComments)

	// carts database CRUD endpoints
	router.GET("/cart/:id", h.GetUserCart)
	router.POST("/cart/add", h.AddProductToCart)
	router.PUT("/cart/:id/update", h.UpdateCart)
	router.DELETE("/cart/:id/delete", h.DeleteCart)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
