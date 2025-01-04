package controller

import (
	"giiku-camp/internal/controller/render"
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/infra/logging"
	"giiku-camp/internal/usecase"
	"giiku-camp/internal/usecase/request"
	_ "giiku-camp/internal/usecase/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCtrl struct {
	UserUsecase usecase.UserUsecase
}

func NewUserCtrl(userUsecase usecase.UserUsecase) UserCtrl {
	return UserCtrl{UserUsecase: userUsecase}
}

// SignUp godoc
// @Summary ユーザー登録
// @Description ユーザー登録
// @Tags User
// @Accept json
// @Produce json
// @Param user body request.SignUpReq true "User details"
// @Success 200 {object} response.SignUpRes "User created"
// @Failure 400 {object} render.Error "Bad request"
// @Router /auth/signup [post]
func (ct *UserCtrl) SignUp(c *gin.Context) {
	req, err := request.NewSignUpReq(c)
	if err != nil {
		logging.Errorf(c, "NewSignUpReq %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ct.UserUsecase.SignUp(c, *req)
	if err != nil {
		switch err {
		case entity.ErrEmailAlreadyUsed:
			logging.Infof(c, "SignUp %v", err)
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeEmailAlreadyUsed)
		default:
			logging.Errorf(c, "SignUp %v", err)
			render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	logging.Infof(c, "SignUp called by %v", res.User)
	render.JSON(c, res)
}

// SignIn godoc
// @Summary ユーザー認証
// @Description ユーザー認証
// @Tags User
// @Accept json
// @Produce json
// @Param user body request.SignInReq true "User details"
// @Success 200 {object} response.SignInRes "User authenticated"
// @Failure 400 {object} render.Error "Bad request"
// @Router /auth/signin [post]
func (ct *UserCtrl) SignIn(c *gin.Context) {
	req, err := request.NewSignInReq(c)
	if err != nil {
		logging.Errorf(c, "NewSignInReq %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ct.UserUsecase.SignIn(c, *req)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			logging.Infof(c, "SignIn %v", err)
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeUserNotFound)
		case entity.ErrPasswordIncorrect:
			logging.Infof(c, "SignIn %v", err)
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodePasswordIncorrect)
		default:
			logging.Errorf(c, "SignIn %v", err)
			render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	logging.Infof(c, "SignIn called by %v", res.User)
	render.JSON(c, res)
}

// RefreshToken godoc
// @Summary トークンの更新
// @Description トークンの更新
// @Tags User
// @Accept json
// @Produce json
// @Param user body request.RefreshTokenReq true "Refresh token"
// @Success 200 {object} response.RefreshTokenRes "Token refreshed"
// @Failure 400 {object} render.Error "Bad request"
// @Router /auth/refresh [post]
func (ct *UserCtrl) RefreshToken(c *gin.Context) {
	req, err := request.NewRefreshTokenReq(c)
	if err != nil {
		logging.Errorf(c, "NewRefreshTokenReq %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ct.UserUsecase.RefreshToken(c, *req)
	if err != nil {
		logging.Infof(c, "RefreshToken %v", err)
		switch err {
		case entity.ErrTokenInValid:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeTokenInValid)
		case entity.ErrFailedToParseClaims:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeFailedToParseClaims)
		case entity.ErrTokenVersionMismatch:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeTokenVersionMismatch)
		default:
			render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	render.JSON(c, res)
}
