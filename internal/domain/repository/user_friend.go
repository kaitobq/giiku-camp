package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"

	"cloud.google.com/go/datastore"
)

type UserFriendRepo interface {
	FindByUserID(ctx context.Context, userId string) ([]entity.UserFriend, error)
	FindByRelatedUserID(ctx context.Context, relatedUserId string) (*entity.UserFriend, error)
	Update(ctx context.Context, userFriendList *entity.UserFriend) error
	UpdateWithTransaction(tx *datastore.Transaction, userFriendList *entity.UserFriend) error
	DeleteWithTransaction(tx *datastore.Transaction, userFriend *entity.UserFriend) error
}
