package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/francoggm/go_expenses_api/configs"
	"github.com/francoggm/go_expenses_api/internal/expenses"
	"github.com/francoggm/go_expenses_api/internal/users"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
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

func configureTimeout() {
	cfg := configs.GetConfigs()

	engine.Use(timeout.New(
		timeout.WithTimeout(cfg.Timeout*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.AbortWithStatus(http.StatusRequestTimeout)
		}),
	))
}

func ConfigureRouters(userHandler *users.Handler, expenseHandler *expenses.Handler) {
	engine = gin.Default()

	configureCors()
	configureTimeout()

	authGroup := engine.Group("/auth")

	authGroup.POST("/signup", userHandler.Register)
	authGroup.POST("/login", userHandler.Login)
	authGroup.POST("/refresh", userHandler.Authenticate, userHandler.RefreshSession)

	expensesGroup := engine.Group("/expenses")

	expensesGroup.Use(userHandler.Authenticate)
	expensesGroup.POST("", expenseHandler.CreateExpense)
	expensesGroup.GET("", expenseHandler.ListExpenses)
	expensesGroup.GET("/:id", expenseHandler.GetExpense)
	expensesGroup.PUT("/:id", expenseHandler.UpdateExpense)
	expensesGroup.DELETE("/:id", expenseHandler.DeleteExpense)
}

func Start(logger *zap.SugaredLogger) error {
	cfg := configs.GetConfigs()

	logger.Infow("starting server",
		zap.String("addr", fmt.Sprintf("%s:%s", cfg.APIAddr, cfg.APIPort)),
	)

	return engine.Run(fmt.Sprintf("%s:%s", cfg.APIAddr, cfg.APIPort))
}
