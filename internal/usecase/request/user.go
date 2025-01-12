package request

import "github.com/gin-gonic/gin"

type SignUpReq struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	GitHubID string `json:"github_id" binding:"omitempty"`
	QiitaID  string `json:"qiita_id" binding:"omitempty"`
	ZennID   string `json:"zenn_id" binding:"omitempty"`
	XID      string `json:"x_id" binding:"omitempty"`
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

type UpdateMeReq struct {
	Name     string `json:"name" binding:"omitempty"`
	GitHubID string `json:"github_id" binding:"omitempty"`
	QiitaID  string `json:"qiita_id" binding:"omitempty"`
	ZennID   string `json:"zenn_id" binding:"omitempty"`
	XID      string `json:"x_id" binding:"omitempty"`
}

func NewUpdateMeReq(c *gin.Context) (*UpdateMeReq, error) {
	var req UpdateMeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
