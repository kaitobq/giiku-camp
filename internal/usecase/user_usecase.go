package usecase

import (
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"
	"giiku-camp/internal/infra/logging"
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
			logging.Errorf(c, "FindByEmail %v", err)
			return nil, err
		}
	}
	if existingUser != nil {
		logging.Infof(c, "FindByEmail returned ErrEmailAlreadyUsed(user: %v)", existingUser.HidePassword())
		return nil, entity.ErrEmailAlreadyUsed
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		logging.Errorf(c, "Create %v", err)
		return nil, err
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID)
	if err != nil {
		logging.Errorf(c, "GenerateAccessToken %v", err)
		return nil, err
	}
	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		logging.Errorf(c, "GenerateRefreshToken %v", err)
		return nil, err
	}

	return response.NewSignUpRes(user, accessToken, refreshToken)
}

func (u *UserUsecase) SignIn(c *gin.Context, req request.SignInReq) (*response.SignInRes, error) {
	ctx := c.Request.Context()
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		logging.Errorf(c, "FindByEmail %v", err)
		return nil, err
	}

	if err := user.VerifyPassword(req.Password); err != nil {
		switch err {
		case entity.ErrPasswordIncorrect:
			logging.Infof(c, "VerifyPassword returned ErrPasswordIncorrect(user: %v)", user.HidePassword())
			return nil, entity.ErrPasswordIncorrect
		default:
			logging.Errorf(c, "VerifyPassword %v", err)
			return nil, err
		}
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID)
	if err != nil {
		logging.Errorf(c, "GenerateAccessToken %v", err)
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		logging.Errorf(c, "GenerateRefreshToken %v", err)
		return nil, err
	}

	return response.NewSignInRes(user, accessToken, refreshToken)
}

func (u *UserUsecase) RefreshToken(c *gin.Context, req request.RefreshTokenReq) (*response.RefreshTokenRes, error) {
	accessToken, refreshToken, err := jwt.RefreshTokens(req.RefreshToken)
	if err != nil {
		logging.Errorf(c, "RefreshTokens %v", err)
		return nil, err
	}

	return response.NewRefreshTokenRes(accessToken, refreshToken)
}
