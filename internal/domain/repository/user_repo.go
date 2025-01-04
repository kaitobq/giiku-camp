package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
)

type UserRepo interface {
	// GetByID(id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	// Update(user *entity.User) error
	// Delete(id uint) error
}
