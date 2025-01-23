package controller

import (
	"giiku-camp/internal/middleware"
	"os"

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

func SetUpRoutes(r *gin.Engine, userCtrl UserCtrl, userCrossingCtrl UserCrossingCtrl, userFriendListCtrl UserFriendListCtrl, middleware *middleware.Middleware) {
	var name, password string
	if name = os.Getenv("BASIC_AUTH_USERNAME"); name == "" {
		panic("BASIC_AUTH_USERNAME is not set")
	}
	if password = os.Getenv("BASIC_AUTH_PASSWORD"); password == "" {
		panic("BASIC_AUTH_PASSWORD is not set")
	}
	r.GET("/swagger/*any", gin.BasicAuth(gin.Accounts{name: password}), ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	v1.GET("/ping", Ping)

	auth := v1.Group("/auth")
	{
		auth.POST("/signup", userCtrl.SignUp)
		auth.POST("/signin", userCtrl.SignIn)
		auth.POST("/refresh", userCtrl.RefreshToken)
	}

	user := v1.Group("/user")
	{
		user.GET("/:id", userCtrl.GetUser)
	}

	authenticated := v1.Group("/authenticated")
	authenticated.Use(middleware.API.Authenticate())
	authenticated.GET("/ping", Ping)

	user = authenticated.Group("/user")
	{
		user.GET("", userCtrl.GetMe)
		user.PUT("", userCtrl.UpdateMe)
	}

	crossing := authenticated.Group("/crossing")
	{
		crossing.POST("", userCrossingCtrl.RegisterUserCrossing)
		crossing.GET("", userCrossingCtrl.GetUserCrossing)
	}

	friendList := authenticated.Group("/friend")
	{
		friendList.GET("", userFriendListCtrl.GetFriendList)
		friendList.POST("", userFriendListCtrl.SendRequest)
		friendList.POST("/accept", userFriendListCtrl.AcceptRequest)
	}
}
