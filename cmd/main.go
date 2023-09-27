package main

import (
	"expenses_api/configs"
	"expenses_api/db"
	"expenses_api/internal/user"
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

	repo := user.NewRepository(db)
	service := user.NewService(repo)
	handler := user.NewHandler(service)
}
