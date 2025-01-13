package repository

import (
	"context"
	"giiku-camp/internal/domain/entity"
)

// MEMO: このinterfaceをinternal/infra/datastore/user_crossing_repo_impl.goで継承する
// TODO: すれ違いAPIを実装するのに必要なDB処理を定義する
type UserCrossingRepo interface {
	Update(ctx context.Context, userCrossing entity.UserCrossing) error
}
