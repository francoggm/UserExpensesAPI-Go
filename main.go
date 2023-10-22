package main

import (
	"log"

	"github.com/francoggm/go_expenses_api/configs"
	"github.com/francoggm/go_expenses_api/db"
	"github.com/francoggm/go_expenses_api/internal/expenses"
	"github.com/francoggm/go_expenses_api/internal/users"
	"github.com/francoggm/go_expenses_api/logger"
	"github.com/francoggm/go_expenses_api/routers"
	"go.uber.org/zap"
)

func main() {
	err := configs.Load()
	if err != nil {
		log.Fatalf("failed to load configs! -> %s", err)
	}

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to create logger! -> %s", err)
	}
	defer logger.Sync()

	db, err := db.NewDatabase()
	if err != nil {
		logger.Fatalw("failed to connect database!",
			zap.Error(err),
		)
	}
	defer db.Close()

	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService, logger)

	expenseRepo := expenses.NewRepository(db)
	expenseService := expenses.NewService(expenseRepo)
	expenseHandler := expenses.NewHandler(expenseService, logger)

	routers.ConfigureRouters(userHandler, expenseHandler)

	err = routers.Start(logger)
	if err != nil {
		logger.Fatalw("failed to start server!",
			zap.Error(err),
		)
	}
}
