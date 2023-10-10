package routers

import (
	"expenses_api/configs"
	"expenses_api/internal/expenses"
	"expenses_api/internal/users"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var engine *gin.Engine

func configureCors() {
	engine.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:8000"},
			AllowMethods:     []string{"GET", "POST", "DELETE", "PUT"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           24 * time.Hour,
		}),
	)
}

func ConfigureRouters(userHandler *users.Handler, expenseHandler *expenses.Handler) {
	engine = gin.Default()

	configureCors()

	authGroup := engine.Group("/auth")

	authGroup.POST("/signup", userHandler.Register)
	authGroup.POST("/login", userHandler.Login)

	expensesGroup := engine.Group("/expenses")

	expensesGroup.Use(userHandler.Authenticate)
	expensesGroup.POST("", expenseHandler.CreateExpense)
	expensesGroup.GET("", expenseHandler.ListExpenses)
	expensesGroup.GET("/:id", expenseHandler.GetExpense)
	expensesGroup.PUT("/:id", expenseHandler.UpdateExpense)
	expensesGroup.DELETE("/:id", expenseHandler.DeleteExpense)
}

func Start() error {
	cfg := configs.GetConfigs()

	return engine.Run(fmt.Sprintf("%s:%s", cfg.APIAddr, cfg.APIPort))
}
