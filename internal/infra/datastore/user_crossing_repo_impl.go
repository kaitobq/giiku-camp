package repository

import (
	"giiku-camp/internal/domain/repository"

	"cloud.google.com/go/datastore"
)

type userCrossingRepo struct {
	db *datastore.Client
}

func NewUserCrossingRepo(db *datastore.Client) repository.UserCrossingRepo {
	return &userCrossingRepo{db: db}
}
