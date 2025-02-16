package request

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserIDNotFound = errors.New("user_id not found")
)

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
	Type   string `json:"type" binding:"required"`
	UserID string // param
}

func NewAcceptRequestReq(c *gin.Context) (*AcceptRequestReq, error) {
	userID := c.Param("user_id")
	var req AcceptRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	req.UserID = userID
	if userID == "" {
		return nil, ErrUserIDNotFound
	}
	return &req, nil
}

type RejectRequestReq struct {
	UserID string // param
}

func NewRejectRequestReq(c *gin.Context) (*RejectRequestReq, error) {
	userID := c.Param("user_id")
	if userID == "" {
		return nil, ErrUserIDNotFound
	}
	return &RejectRequestReq{UserID: userID}, nil
}
