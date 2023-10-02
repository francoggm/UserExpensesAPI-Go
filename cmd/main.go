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
		log.Fatalf("failed to load configs! -> %s", err)
	}

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect database!-> %s", err)
	}
	defer db.Close()

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	routers.ConfigureRouters(userHandler)
	
	err = routers.Start()
	if err != nil {
		log.Fatal("failed to start server!")
	}
}
