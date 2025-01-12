package controller

import (
	"giiku-camp/internal/middleware"

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

func SetUpRoutes(r *gin.Engine, userCtrl UserCtrl, userCrossingCtrl UserCrossingCtrl, middleware *middleware.Middleware) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	v1.GET("/ping", Ping)

	auth := v1.Group("/auth")
	{
		auth.POST("/signup", userCtrl.SignUp)
		auth.POST("/signin", userCtrl.SignIn)
		auth.POST("/refresh", userCtrl.RefreshToken)
	}

	authenticated := v1.Group("/authenticated")
	authenticated.Use(middleware.API.Authenticate())
	authenticated.GET("/ping", Ping)

	user := authenticated.Group("/user")
	{
		user.GET("", userCtrl.GetMe)
		user.PUT("", userCtrl.UpdateMe)
	}

	_ = v1.Group("/crossing")
	{
		// TODO: すれ違いAPIのエンドポイントを追加する
	}
}
