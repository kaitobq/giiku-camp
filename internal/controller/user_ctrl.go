package controller

import (
	"giiku-camp/internal/controller/render"
	"giiku-camp/internal/usecase"
	"giiku-camp/internal/usecase/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCtrl struct {
	UserUsecase usecase.UserUsecase
}

func NewUserCtrl(userUsecase usecase.UserUsecase) UserCtrl {
	return UserCtrl{UserUsecase: userUsecase}
}

func (ct *UserCtrl) CreateUser(c *gin.Context) {
	req, err := request.NewCreateUserRequest(c)
	if err != nil {
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
	}

	res, err := ct.UserUsecase.CreateUser(c, *req)
	if err != nil {
		render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
	}

	render.JSON(c, res)
}
