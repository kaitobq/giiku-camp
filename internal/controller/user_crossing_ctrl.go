package controller

import (
	"giiku-camp/internal/controller/render"
	"giiku-camp/internal/infra/logging"
	"giiku-camp/internal/usecase"
	"giiku-camp/internal/usecase/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCrossingCtrl struct {
	UserCrossingUsecase usecase.UserCrossingUsecase
}

func NewUserCrossingCtrl(userCrossingUsecase usecase.UserCrossingUsecase) UserCrossingCtrl {
	return UserCrossingCtrl{UserCrossingUsecase: userCrossingUsecase}
}

// TODO: すれ違いAPIを実装する
func (ct *UserCrossingCtrl) RegisterUserCrossing(c *gin.Context) {
	req, err := request.NewRegisterUserCrossingReq(c)
	if err != nil {
		logging.Errorf(c, "NewRegisterUserCrossingReq %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ct.UserCrossingUsecase.RegisterUserCrossing(c, *req)
	if err != nil {
		logging.Errorf(c, "RegisterUserCrossing %v", err)
		return
	}

	render.JSON(c, res)
}
