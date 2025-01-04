package request

import "github.com/gin-gonic/gin"

type SignUpReq struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewSignUpReq(c *gin.Context) (*SignUpReq, error) {
	var req SignUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewSignInReq(c *gin.Context) (*SignInReq, error) {
	var req SignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
