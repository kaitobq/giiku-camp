package response

import (
	"giiku-camp/internal/domain/entity"
	"time"
)

type UserRes struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserRes(user *entity.User) (*UserRes, error) {
	return &UserRes{
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

type TokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpRes struct {
	User  UserRes  `json:"user"`
	Token TokenRes `json:"token"`
}

func NewSignUpRes(user *entity.User, accessToken, refreshToken string) (*SignUpRes, error) {
	res := SignUpRes{}
	userRes, err := NewUserRes(user)
	if err != nil {
		return nil, err
	}
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
	userRes, err := NewUserRes(user)
	if err != nil {
		return nil, err
	}
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
