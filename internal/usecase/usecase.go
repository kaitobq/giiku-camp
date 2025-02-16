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
	GetUser(c *gin.Context, userID string) (*response.CrossedUserRes, error)
	UpdateMe(c *gin.Context, req request.UpdateMeReq) (*response.UserRes, error)
}

type UserCrossingUsecase interface {
	RegisterUserCrossing(c *gin.Context, req request.RegisterUserCrossingReq) (*response.RegisterUserCrossingRes, error)
	GetUserCrossing(c *gin.Context) (*response.GetUserCrossingRes, error)
}

type UserFriendUsecase interface {
	GetUserFriend(c *gin.Context) (*response.UserFriendRes, error)
	SendRequest(c *gin.Context, req request.SendRequestReq) (*response.SendRequestRes, error)
	AcceptRequest(c *gin.Context, req request.AcceptRequestReq) (*response.AcceptRequestRes, error)
	RejectRequest(c *gin.Context, req request.RejectRequestReq) (*response.RejectRequestRes, error)
}
