package usecase

import "giiku-camp/internal/domain/repository"

type userCrossingUsecase struct {
	userCrossingRepo repository.UserCrossingRepo
}

func NewUserCrossingUsecase(userCrossingRepo repository.UserCrossingRepo) UserCrossingUsecase {
	return &userCrossingUsecase{userCrossingRepo: userCrossingRepo}
}
