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

// CreateUser godoc
// @Summary Create a user
// @Description Create new user by providing user details
// @Tags User
// @Accept json
// @Produce json
// @Param user body map[string]interface{} true "User details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users [post]
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
