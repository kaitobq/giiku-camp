package repository

import (
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
