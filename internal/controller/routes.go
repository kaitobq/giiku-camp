package controller

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Ping godoc
// @Summary Ping the server
// @Description Returns "pong" if server is working
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func SetUpRoutes(r *gin.Engine, userCtrl UserCtrl) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", Ping)

	auth := r.Group("/auth")
	{
		auth.POST("/signup", userCtrl.SignUp)
	}
}
