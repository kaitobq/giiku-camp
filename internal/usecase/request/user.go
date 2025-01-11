package request

import "github.com/gin-gonic/gin"

type SignUpReq struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewSignUpReq(c *gin.Context) (*SignUpReq, error) {
	var req SignUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type SignInReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewSignInReq(c *gin.Context) (*SignInReq, error) {
	var req SignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func NewRefreshTokenReq(c *gin.Context) (*RefreshTokenReq, error) {
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
