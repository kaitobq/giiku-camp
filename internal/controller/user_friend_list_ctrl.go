package controller

import (
	"giiku-camp/internal/controller/render"
	"giiku-camp/internal/domain/xcontext"
	"giiku-camp/internal/usecase"
	_ "giiku-camp/internal/usecase/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserFriendListCtrl struct {
	UserFriendListUseCase usecase.UserFriendListUsecase
}

func NewUserFriendListCtrl(userFriendListUseCase usecase.UserFriendListUsecase) UserFriendListCtrl {
	return UserFriendListCtrl{UserFriendListUseCase: userFriendListUseCase}
}

// GetFriendList godoc
// @Summary フレンド情報取得
// @Description フレンド情報取得
// @Tags Friend
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserFriendListRes "Friend details"
// @Failure 400 {object} render.Error "Bad request"
// @Failure 500 {object} render.Error "Internal server error"
// @Router /api/v1/authenticated/friend [get]
func (ct *UserFriendListCtrl) GetFriendList(c *gin.Context) {
	user := xcontext.User(c)

	res, err := ct.UserFriendListUseCase.GetUserFriendList(c, user.ID)
	if err != nil {
		render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(c, res)
}
