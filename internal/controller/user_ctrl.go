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
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeEmailAlreadyUsed)
		default:
			logging.Errorf(c, "SignUp %v", err)
			render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	render.JSON(c, res)
}
