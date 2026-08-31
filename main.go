package main

import (
	"ecommerce/database"
	"ecommerce/middleware"
	"ecommerce/routers"
	"log"
)

func main() {
	// Initialize database connection and run schema migrations
	sqlDB, db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the database connection is closed when the application terminates
	defer sqlDB.Close()

	// Initialize router and register API endpoints
	router := routers.SetupRouters(db)

	// Apply global security middleware
	middleware.ApplySecurityMiddlewares(router)

	// Start the HTTP server
	router.Run()
}
