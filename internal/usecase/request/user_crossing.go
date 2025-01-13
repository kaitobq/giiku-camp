package request

import "github.com/gin-gonic/gin"

// TODO: すれ違いAPIのリクエストを定義する
type RegisterUserCrossingReq struct {
	UserIDs []string `json:"user_ids" binding:"required"`
}

func NewRegisterUserCrossingReq(c *gin.Context) (*RegisterUserCrossingReq, error) {
	var req RegisterUserCrossingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
