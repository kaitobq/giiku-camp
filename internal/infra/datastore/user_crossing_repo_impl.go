package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"

	"cloud.google.com/go/datastore"
)

// TODO: repository.UserCrossingRepoを継承する
type userCrossingRepo struct {
	db *datastore.Client
}

func NewUserCrossingRepo(db *datastore.Client) repository.UserCrossingRepo {
	return &userCrossingRepo{db: db}
}

func (r *userCrossingRepo) Update(ctx context.Context, userCrossing entity.UserCrossing) error {
	k := datastore.NameKey("UserCrossing", userCrossing.ID, nil)
	userCrossing.UpdateUpdatedAt()
	if _, err := r.db.Put(ctx, k, userCrossing); err != nil {
		return err
	}
	return nil
}
