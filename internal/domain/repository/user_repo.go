package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"

	"cloud.google.com/go/datastore"
)

type UserRepo interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	UpdateWithTransaction(tx *datastore.Transaction, user *entity.User) error
	FindByAppleID(ctx context.Context, appleID string) (*entity.User, error)
}
