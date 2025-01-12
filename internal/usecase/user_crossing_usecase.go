package usecase

import "giiku-camp/internal/domain/repository"

// TODO: UserCrossingUsecaseを継承する
type userCrossingUsecase struct {
	userCrossingRepo repository.UserCrossingRepo
}

func NewUserCrossingUsecase(userCrossingRepo repository.UserCrossingRepo) UserCrossingUsecase {
	return &userCrossingUsecase{userCrossingRepo: userCrossingRepo}
}
