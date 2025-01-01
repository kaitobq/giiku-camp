package usecase

import (
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"

	"github.com/gin-gonic/gin"
)

type UserUsecase struct {
	userRepo repository.UserRepo
}

func NewUserUsecase(userRepo repository.UserRepo) UserUsecase {
	return UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) CreateUser(c *gin.Context, req request.CreateUserRequest) (*response.CreateUserResponse, error) {
	user, err := entity.NewUser(req.Email, req.UserName, req.Password)
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return response.NewCreateUserResponse(user)
}
