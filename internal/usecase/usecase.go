package usecase

import (
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"

	"github.com/gin-gonic/gin"
)

type UserUsecase interface {
	SignUp(c *gin.Context, req request.SignUpReq) (*response.SignUpRes, error)
	SignIn(c *gin.Context, req request.SignInReq) (*response.SignInRes, error)
	RefreshToken(c *gin.Context, req request.RefreshTokenReq) (*response.RefreshTokenRes, error)
	GetMe(c *gin.Context) (*response.UserRes, error)
	UpdateMe(c *gin.Context, req request.UpdateMeReq) (*response.UserRes, error)
}

type UserCrossingUsecase interface {
	RegisterUserCrossing(c *gin.Context, req request.RegisterUserCrossingReq) (*response.RegisterUserCrossingRes, error)
	GetUserCrossing(c *gin.Context) (*response.GetUserCrossingRes, error)
}
