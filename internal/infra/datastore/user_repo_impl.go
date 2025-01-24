package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type userRepo struct {
	db *datastore.Client
}

func NewUserRepo(db *datastore.Client) repository.UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	k := datastore.NameKey("User", id, nil)
	var user entity.User
	if err := r.db.Get(ctx, k, &user); err != nil {
		// if err == datastore.ErrNoSuchEntity {
		// 	return nil, entity.ErrUserNotFound
		// }
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	k := datastore.NameKey("User", user.ID, nil)
	user.UpdateUpdatedAt()
	if _, err := r.db.Put(ctx, k, user); err != nil {
		return err
	}
	return nil
}

func (r *userRepo) FindByAppleID(ctx context.Context, appleID string) (*entity.User, error) {
	q := datastore.NewQuery("User").FilterField("AppleID", "=", appleID).Limit(1)
	it := r.db.Run(ctx, q)
	var user entity.User
	_, err := it.Next(&user)
	if err == iterator.Done {
		return nil, entity.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
