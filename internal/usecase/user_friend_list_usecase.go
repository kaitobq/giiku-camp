package usecase

import "giiku-camp/internal/domain/repository"

type userFriendListUsecase struct {
	userFriendListRepo repository.UserFriendListRepo
}

func NewUserFriendListUsecase(userFriendListRepo repository.UserFriendListRepo) UserFriendListUsecase {
	return &userFriendListUsecase{userFriendListRepo: userFriendListRepo}
}
