package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"

	"cloud.google.com/go/datastore"
)

type UserFriendListRepo interface {
	FindByUserID(ctx context.Context, userId string) (*entity.UserFriendList, error)
	Update(ctx context.Context, userFriendList *entity.UserFriendList) error
	UpdateWithTransaction(tx *datastore.Transaction, userFriendList *entity.UserFriendList) error
}
