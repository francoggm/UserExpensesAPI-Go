package main

import (
	"expenses_api/configs"
	"expenses_api/db"
	"expenses_api/internal/users"
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

	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	routers.ConfigureRouters(userHandler)

	err = routers.Start()
	if err != nil {
		log.Fatalf("failed to start server!-> %s", err)
	}
}
