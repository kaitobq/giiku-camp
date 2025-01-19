package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
)

type UserFriendListRepo interface {
	FindByUserID(ctx context.Context, userId string) ([]entity.UserFriendList, error)
	Update(ctx context.Context, userFriendList *entity.UserFriendList) error
}
