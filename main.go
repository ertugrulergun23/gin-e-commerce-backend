package main

import (
	"ecommerce/database"
	"ecommerce/routers"
	"log"
)

func main() {
	sqlDB, db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()
	router := routers.SetupRouters(db)
	router.Run()
}
