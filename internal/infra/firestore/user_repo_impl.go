package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	iter := r.db.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		// ドキュメントが見つからなかった場合
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
