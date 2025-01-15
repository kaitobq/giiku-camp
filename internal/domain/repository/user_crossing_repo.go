package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
)

type UserCrossingRepo interface {
	FindByUserID(ctx context.Context, userId string) ([]string, error)
	Update(ctx context.Context, userCrossing entity.UserCrossing) error
}
