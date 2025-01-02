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
	user, err := entity.NewUser(req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	existingUser, err := u.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			// ユーザーが見つからなかった場合はそのまま処理を続ける
		default:
			return nil, err
		}
	}
	if existingUser != nil {
		return nil, entity.ErrEmailAlreadyUsed
	}

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
