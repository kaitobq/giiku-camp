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
	res, err := ct.UserFriendListUseCase.GetUserFriendList(c)
	if err != nil {
		logging.Errorf(c, "GetUserFriendList %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(c, res)
}

// SendRequest godoc
// @Summary フレンド依頼送信
// @Description フレンド依頼送信
// @Tags Friend
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body request.SendRequestReq true "Friend details"
// @Success 200 {object} response.SendRequestRes "Friend details"
// @Failure 400 {object} render.Error "Bad request"
// @Failure 500 {object} render.Error "Internal server error"
// @Router /api/v1/authenticated/friend [post]
func (ct *UserFriendListCtrl) SendRequest(c *gin.Context) {
	req, err := request.NewSendRequestReq(c)
	if err != nil {
		logging.Errorf(c, "NewSendRequestReq %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ct.UserFriendListUseCase.SendRequest(c, *req)
	if err != nil {
		logging.Errorf(c, "SendRequest %v", err)
		switch err {
		case entity.ErrFriendRequestAlreadySent:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeFriendRequestAlreadySent)
		case entity.ErrFriendRequestAlreadyReceived:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeFriendRequestAlreadyReceived)
		case entity.ErrAlreadyFriend:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeAlreadyFriend)
		default:
			render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	render.JSON(c, res)
}

// AcceptRequest godoc
// @Summary フレンド依頼許可
// @Description フレンド依頼許可
// @Tags Friend
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body request.AcceptRequestReq true "Friend details"
// @Success 200 {object} response.AcceptRequestRes "Friend details"
// @Failure 400 {object} render.Error "Bad request"
// @Failure 500 {object} render.Error "Internal server error"
// @Router /api/v1/authenticated/friend/accept [post]
func (ct *UserFriendListCtrl) AcceptRequest(c *gin.Context) {
	req, err := request.NewAcceptRequestReq(c)
	if err != nil {
		logging.Errorf(c, "NewAcceptRequestReq %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ct.UserFriendListUseCase.AcceptRequest(c, *req)
	if err != nil {
		logging.Errorf(c, "AcceptRequest %v", err)
		switch err {
		case entity.ErrFriendRequestNotFound:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeFriendRequestNotFound)
		case entity.ErrSentRequestNotFound:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeSentRequestNotFound)
		default:
			render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	render.JSON(c, res)
}
