package controller

import "giiku-camp/internal/usecase"

type UserFriendListCtrl struct {
	UserFriendListUseCase usecase.UserFriendListUsecase
}

func NewUserFriendListCtrl(userFriendListUseCase usecase.UserFriendListUsecase) UserFriendListCtrl {
	return UserFriendListCtrl{UserFriendListUseCase: userFriendListUseCase}
}
