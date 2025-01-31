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

type AcceptRequestReq struct {
	UserID string `json:"user_id" binding:"required"`
}

func NewAcceptRequestReq(c *gin.Context) (*AcceptRequestReq, error) {
	var req AcceptRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type RejectRequestReq struct {
	UserID string `json:"user_id" binding:"required"`
}

func NewRejectRequestReq(c *gin.Context) (*RejectRequestReq, error) {
	var req RejectRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
