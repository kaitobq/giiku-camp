package response

import (
	"giiku-camp/internal/domain/entity"
	"time"
)

type UserRes struct {
	Email     string    `json:"email"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserRes(user *entity.User) (*UserRes, error) {
	return &UserRes{
		Email:     user.Email,
		UserName:  user.UserName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

type SignUpRes struct {
	User  UserRes `json:"user"`
	Token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"token"`
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
