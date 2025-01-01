package response

import (
	"giiku-camp/internal/domain/entity"
	"time"
)

type CreateUserResponse struct {
	Email     string    `json:"email"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCreateUserResponse(user *entity.User) (*CreateUserResponse, error) {
	return &CreateUserResponse{
		Email:     user.Email,
		UserName:  user.UserName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
