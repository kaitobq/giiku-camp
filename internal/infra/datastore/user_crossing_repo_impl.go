package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"

	"cloud.google.com/go/datastore"
)

type userCrossingRepo struct {
	db *datastore.Client
}

func NewUserCrossingRepo(db *datastore.Client) repository.UserCrossingRepo {
	return &userCrossingRepo{db: db}
}

func (r *userCrossingRepo) Update(ctx context.Context, userCrossing entity.UserCrossing) error {
	k := datastore.NameKey("UserCrossing", userCrossing.ID, nil)
	userCrossing.UpdateUpdatedAt()
	if _, err := r.db.Put(ctx, k, &userCrossing); err != nil {
		return err
	}
	return nil
}

func (r *userCrossingRepo) FindByUserID(ctx context.Context, userId string) ([]*entity.UserCrossing, error) {
	q := datastore.NewQuery("UserCrossing").FilterField("userid", "=", userId)
	var userCrossings []*entity.UserCrossing
	keys, err := r.db.GetAll(ctx, q, &userCrossings)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		userCrossings[i].ID = key.Name
	}
	return userCrossings, nil
}
