package controller

import "giiku-camp/internal/usecase"

type UserCrossingCtrl struct {
	UserCrossingUsecase usecase.UserCrossingUsecase
}

func NewUserCrossingCtrl(userCrossingUsecase usecase.UserCrossingUsecase) UserCrossingCtrl {
	return UserCrossingCtrl{UserCrossingUsecase: userCrossingUsecase}
}

// TODO: すれ違いAPIを実装する
