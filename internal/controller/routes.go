package controller

import "github.com/gin-gonic/gin"

func SetUpRoutes(r *gin.Engine, userCtrl UserCtrl) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/users", userCtrl.CreateUser)
}
