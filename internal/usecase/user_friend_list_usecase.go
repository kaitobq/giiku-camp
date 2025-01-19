package usecase

import (
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"
	"giiku-camp/internal/usecase/response"

	"github.com/gin-gonic/gin"
)

type userFriendListUsecase struct {
	userFriendListRepo repository.UserFriendListRepo
}

func NewUserFriendListUsecase(userFriendListRepo repository.UserFriendListRepo) UserFriendListUsecase {
	return &userFriendListUsecase{userFriendListRepo: userFriendListRepo}
}

func (u *userFriendListUsecase) GetUserFriendList(c *gin.Context, userID string) (*response.UserFriendListRes, error) {
	userFriendList, err := u.userFriendListRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		if err == entity.ErrUserFriendListNotFound {
			ent := entity.NewUserFriendList(userID)
			if err := u.userFriendListRepo.Update(c.Request.Context(), ent); err != nil {
				return nil, err
			}
			return response.NewUserFriendListRes(*ent), nil
		}
		return nil, err
	}
	return response.NewUserFriendListRes(*userFriendList), nil
}
