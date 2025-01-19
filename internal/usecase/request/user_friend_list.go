package request

import "github.com/gin-gonic/gin"

type SendRequestReq struct {
	UserID string `json:"user_id" binding:"required"`
}

func NewSendRequestReq(c *gin.Context) (*SendRequestReq, error) {
	var req SendRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
