package routers

import (
	"expenses_api/configs"
	"expenses_api/internal/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

func configureCors(engine *gin.Engine){

}

func configureMiddlewares(engine *gin.Engine){

}

func ConfigureRouters(engine *gin.Engine, userHandler *user.Handler) {
	configureCors(engine)
	configureMiddlewares(engine)

	authGroup := engine.Group("/auth")
	authGroup.POST("/signup", userHandler.Register)	
	authGroup.POST("/login", userHandler.Login)	
}

func Start(engine *gin.Engine) error {
	cfg := configs.GetConfigs()

	return engine.Run(fmt.Sprintf("%s:%s", cfg.APIAddr, cfg.APIPort)) 
}