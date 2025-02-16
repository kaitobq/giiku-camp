package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"

	"cloud.google.com/go/datastore"
)

type userFriendRepo struct {
	db *datastore.Client
}

func NewUserFriendRepo(db *datastore.Client) repository.UserFriendRepo {
	return &userFriendRepo{db: db}
}

func (r *userFriendRepo) FindByUserID(ctx context.Context, userId string) ([]entity.UserFriend, error) {
	q := datastore.NewQuery("UserFriend").FilterField("UserID", "=", userId)
	var userFriends []entity.UserFriend
	_, err := r.db.GetAll(ctx, q, &userFriends)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, entity.ErrUserFriendNotFound
		}
		return nil, err
	}
	return userFriends, nil
}

func (r *userFriendRepo) FindByRelatedUserID(ctx context.Context, relatedUserId string) (*entity.UserFriend, error) {
	k := datastore.NameKey("UserFriend", relatedUserId, nil)
	var userFriend entity.UserFriend
	if err := r.db.Get(ctx, k, &userFriend); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, entity.ErrUserFriendNotFound
		}
		return nil, err
	}
	return &userFriend, nil
}

func (r *userFriendRepo) Update(ctx context.Context, userFriend *entity.UserFriend) error {
	k := datastore.NameKey("UserFriend", userFriend.UserID, nil)
	if _, err := r.db.Put(ctx, k, userFriend); err != nil {
		return err
	}
	return nil
}

func (r *userFriendRepo) UpdateWithTransaction(tx *datastore.Transaction, userFriend *entity.UserFriend) error {
	k := datastore.NameKey("UserFriend", userFriend.UserID, nil)
	if _, err := tx.Put(k, userFriend); err != nil {
		return err
	}
	return nil
}

func (r *userFriendRepo) DeleteWithTransaction(tx *datastore.Transaction, userFriend *entity.UserFriend) error {
	k := datastore.NameKey("UserFriend", userFriend.UserID, nil)
	if err := tx.Delete(k); err != nil {
		return err
	}
	return nil
}
