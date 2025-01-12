package request

import "github.com/gin-gonic/gin"

type SignUpReq struct {
	Email    string `json:"email" binding:"required" example:"someone@example.com"`
	Name     string `json:"name" binding:"required" example:"someone"`
	Password string `json:"password" binding:"required" example:"password123"`
	GitHubID string `json:"github_id" binding:"omitempty" example:"someone"`
	QiitaID  string `json:"qiita_id" binding:"omitempty" example:"someone"`
	ZennID   string `json:"zenn_id" binding:"omitempty" example:"someone"`
	XID      string `json:"x_id" binding:"omitempty" example:"someone"`
}

func NewSignUpReq(c *gin.Context) (*SignUpReq, error) {
	var req SignUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type SignInReq struct {
	Email    string `json:"email" binding:"required" example:"someone@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

func NewSignInReq(c *gin.Context) (*SignInReq, error) {
	var req SignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkyMzIwMzMsImlhdCI6MTczNjY0MDAzMywidG9rZW5fdmVyc2lvbiI6NiwidXNlcl9pZCI6IjNlZTAzYWRkLWJjZDItNDU1Zi05OTk0LWY3YWY4NDZjZTY2MiJ9.WJCX_65UJKX-DjUpr9TvbtXsE6ZyUH6NyCLNGwcMWR0"`
}

func NewRefreshTokenReq(c *gin.Context) (*RefreshTokenReq, error) {
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

type UpdateMeReq struct {
	Name     string `json:"name" binding:"omitempty" example:"someone"`
	GitHubID string `json:"github_id" binding:"omitempty" example:"someone"`
	QiitaID  string `json:"qiita_id" binding:"omitempty" example:"someone"`
	ZennID   string `json:"zenn_id" binding:"omitempty" example:"someone"`
	XID      string `json:"x_id" binding:"omitempty" example:"someone"`
}

func NewUpdateMeReq(c *gin.Context) (*UpdateMeReq, error) {
	var req UpdateMeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
