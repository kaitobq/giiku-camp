package response

import (
	"giiku-camp/internal/domain/entity"
	"time"
)

type UserRes struct {
	ID        string    `json:"id" example:"3ee03add-bcd2-455f-9994-f7af846ce662"`
	Email     string    `json:"email" example:"someone@example.com"`
	Name      string    `json:"name" example:"someone"`
	GitHubID  string    `json:"github_id" example:"someone"`
	QiitaID   string    `json:"qiita_id" example:"someone"`
	ZennID    string    `json:"zenn_id" example:"someone"`
	XID       string    `json:"x_id" example:"someone"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-10T16:45:40.876773Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-01-12T00:21:57.251469Z"`
}

func NewUserRes(user *entity.User) *UserRes {
	return &UserRes{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		GitHubID:  user.GitHubID,
		QiitaID:   user.QiitaID,
		ZennID:    user.ZennID,
		XID:       user.XID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type TokenRes struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY2NDM2MzMsImlhdCI6MTczNjY0MDAzMywidXNlcl9pZCI6IjNlZTAzYWRkLWJjZDItNDU1Zi05OTk0LWY3YWY4NDZjZTY2MiJ9.QbduICi7TKkVnRckCJbTCYurvmBnXQlmclSm7BKsLxo"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkyMzIwMzMsImlhdCI6MTczNjY0MDAzMywidG9rZW5fdmVyc2lvbiI6NiwidXNlcl9pZCI6IjNlZTAzYWRkLWJjZDItNDU1Zi05OTk0LWY3YWY4NDZjZTY2MiJ9.WJCX_65UJKX-DjUpr9TvbtXsE6ZyUH6NyCLNGwcMWR0"`
}

type SignUpRes struct {
	User  UserRes  `json:"user"`
	Token TokenRes `json:"token"`
}

func NewSignUpRes(user *entity.User, accessToken, refreshToken string) (*SignUpRes, error) {
	res := SignUpRes{}
	userRes := NewUserRes(user)
	res.User = *userRes
	res.Token.AccessToken = accessToken
	res.Token.RefreshToken = refreshToken
	return &res, nil
}

type SignInRes struct {
	User  UserRes  `json:"user"`
	Token TokenRes `json:"token"`
}

func NewSignInRes(user *entity.User, accessToken, refreshToken string) (*SignInRes, error) {
	res := SignInRes{}
	userRes := NewUserRes(user)
	res.User = *userRes
	res.Token.AccessToken = accessToken
	res.Token.RefreshToken = refreshToken
	return &res, nil
}

type RefreshTokenRes struct {
	Token TokenRes `json:"token"`
}

func NewRefreshTokenRes(accessToken, refreshToken string) (*RefreshTokenRes, error) {
	res := RefreshTokenRes{}
	res.Token.AccessToken = accessToken
	res.Token.RefreshToken = refreshToken
	return &res, nil
}
