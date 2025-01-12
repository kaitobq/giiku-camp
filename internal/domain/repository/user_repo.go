package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}
