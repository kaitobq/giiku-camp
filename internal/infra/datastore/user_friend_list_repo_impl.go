package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"

	"cloud.google.com/go/datastore"
)

type userFriendListRepo struct {
	db *datastore.Client
}

func NewUserFriendListRepo(db *datastore.Client) repository.UserFriendListRepo {
	return &userFriendListRepo{db: db}
}

func (r *userFriendListRepo) FindByUserID(ctx context.Context, userId string) (*entity.UserFriendList, error) {
	k := datastore.NameKey("UserFriendList", userId, nil)
	var userFriendList entity.UserFriendList
	if err := r.db.Get(ctx, k, &userFriendList); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, entity.ErrUserFriendListNotFound
		}
		return nil, err
	}
	return &userFriendList, nil
}

func (r *userFriendListRepo) Update(ctx context.Context, userFriendList *entity.UserFriendList) error {
	k := datastore.NameKey("UserFriendList", userFriendList.UserID, nil)
	if _, err := r.db.Put(ctx, k, userFriendList); err != nil {
		return err
	}
	return nil
}
