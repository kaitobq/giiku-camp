package controller

import (
	"giiku-camp/internal/controller/render"
	"giiku-camp/internal/infra/logging"
	"giiku-camp/internal/usecase"
	"giiku-camp/internal/usecase/request"
	_ "giiku-camp/internal/usecase/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCrossingCtrl struct {
	UserCrossingUsecase usecase.UserCrossingUsecase
}

func NewUserCrossingCtrl(userCrossingUsecase usecase.UserCrossingUsecase) UserCrossingCtrl {
	return UserCrossingCtrl{UserCrossingUsecase: userCrossingUsecase}
}

// RegisterUserCrossing godoc
// @Summary すれ違ったユーザの取得
// @Description すれ違ったユーザの取得
// @Tags Crossing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body request.RegisterUserCrossingReq true "User Crossing Details"
// @Success 200 {object} response.RegisterUserCrossingRes "UserCrossing created"
// @Failure 400 {object} render.Error "Bad request"
// @Failure 500 {object} render.Error "Internal server error"
// @Router /api/v1/authenticated/crossing [post]
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
		render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(c, res)
}
