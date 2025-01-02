package usecase

import (
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"
	"giiku-camp/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type UserUsecase struct {
	userRepo repository.UserRepo
}

func NewUserUsecase(userRepo repository.UserRepo) UserUsecase {
	return UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) SignUp(c *gin.Context, req request.SignUpReq) (*response.SignUpRes, error) {
	user, err := entity.NewUser(req.Email, req.Name, req.Password)
	if err != nil {
		return nil, err
	}

	existingUser, err := u.userRepo.FindByEmail(c.Request.Context(), user.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, entity.ErrEmailAlreadyUsed
	}

	ctx := c.Request.Context()
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return response.NewSignUpRes(user, accessToken, refreshToken)
}
