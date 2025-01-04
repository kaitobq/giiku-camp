package xcontext

import (
	"giiku-camp/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

var (
	authUserKey = "AUTH:ME"
)

func WithUser(c *gin.Context, user *entity.User) {
	c.Set(authUserKey, user)
}

func User(c *gin.Context) *entity.User {
	if v, ok := c.Get(authUserKey); ok {
		return v.(*entity.User)
	}
	return nil
}
