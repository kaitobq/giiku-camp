package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"

	"cloud.google.com/go/firestore"
)

type userRepo struct {
	db *firestore.Client
}

func NewUserRepo(db *firestore.Client) repository.UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	user.UpdateCreatedAt()
	user.UpdateUpdatedAt()
	_, err := r.db.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
