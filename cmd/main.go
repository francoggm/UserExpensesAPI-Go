package main

import (
	"expenses_api/configs"
	"expenses_api/db"
	"expenses_api/internal/user"
	"expenses_api/routers"
	"log"
)

func main() {
	err := configs.Load()
	if err != nil {
		log.Fatal("failed to load configs!")
	}

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal("failed to connect database!")
	}
	defer db.Close()

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	routers.InitRouter(userHandler)
}
