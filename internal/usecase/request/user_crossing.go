package request

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrEmptyUserIDs = errors.New("user ids cannot be empty")
)

var (
	CodeEmptyUserIDs = 21000
)

type RegisterUserCrossingReq struct {
	UserIDs []string `json:"user_ids" binding:"required"`
}

func NewRegisterUserCrossingReq(c *gin.Context) (*RegisterUserCrossingReq, error) {
	var req RegisterUserCrossingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	if len(req.UserIDs) == 0 {
		return nil, ErrEmptyUserIDs
	}
	return &req, nil
}
