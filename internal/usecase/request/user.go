package request

import "github.com/gin-gonic/gin"

type CreateUserRequest struct {
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func NewCreateUserRequest(c *gin.Context) (*CreateUserRequest, error) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
