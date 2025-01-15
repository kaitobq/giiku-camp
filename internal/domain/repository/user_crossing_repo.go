package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
)

type UserCrossingRepo interface {
	FindByUserID(ctx context.Context, userId string) ([]entity.UserCrossing, error)
	Update(ctx context.Context, userCrossing entity.UserCrossing) error
}
