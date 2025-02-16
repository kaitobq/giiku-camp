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

type UserFriendCtrl struct {
	UserFriendUsecase usecase.UserFriendUsecase
}

func NewUserFriendCtrl(userFriendUsecase usecase.UserFriendUsecase) UserFriendCtrl {
	return UserFriendCtrl{UserFriendUsecase: userFriendUsecase}
}

// GetUserFriend godoc
// @Summary フレンド情報取得
// @Description フレンド情報取得
// @Tags Friend
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserFriendRes "Friend details"
// @Failure 400 {object} render.Error "Bad request"
// @Failure 500 {object} render.Error "Internal server error"
// @Router /api/v1/authenticated/friend [get]
func (ct *UserFriendCtrl) GetUserFriend(c *gin.Context) {
	res, err := ct.UserFriendUsecase.GetUserFriend(c)
	if err != nil {
		logging.Errorf(c, "GetUserFriend %v", err)
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
func (ct *UserFriendCtrl) SendRequest(c *gin.Context) {
	req, err := request.NewSendRequestReq(c)
	if err != nil {
		logging.Errorf(c, "NewSendRequestReq %v", err)
		render.ErrorJSON(c, err.Error(), 400)
		return
	}

	res, err := ct.UserFriendUsecase.SendRequest(c, *req)
	if err != nil {
		logging.Errorf(c, "SendRequest %v", err)
		switch err {
		case entity.ErrUserFriendNotFound:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeUserFriendNotFound)
		case entity.ErrAlreadyFriend:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeAlreadyFriend)
		case entity.ErrFriendRequestAlreadySent:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeFriendRequestAlreadySent)
		case entity.ErrFriendRequestAlreadyReceived:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeFriendRequestAlreadyReceived)
		default:
			render.ErrorJSON(c, err.Error(), 500)
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
// @Router /api/v1/authenticated/friend/{user_id} [patch]
func (ct *UserFriendCtrl) AcceptRequest(c *gin.Context) {
	req, err := request.NewAcceptRequestReq(c)
	if err != nil {
		logging.Errorf(c, "NewAcceptRequestReq %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ct.UserFriendUsecase.AcceptRequest(c, *req)
	if err != nil {
		logging.Errorf(c, "AcceptRequest %v", err)
		switch err {
		case entity.ErrUserFriendNotFound:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeUserFriendNotFound)
		default:
			render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	render.JSON(c, res)
}

// RejectRequest godoc
// @Summary フレンド依頼拒否
// @Description フレンド依頼拒否
// @Tags Friend
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.RejectRequestRes "Friend details"
// @Failure 400 {object} render.Error "Bad request"
// @Failure 500 {object} render.Error "Internal server error"
// @Router /api/v1/authenticated/friend/{user_id} [delete]
func (ct *UserFriendCtrl) RejectRequest(c *gin.Context) {
	req, err := request.NewRejectRequestReq(c)
	if err != nil {
		logging.Errorf(c, "NewRejectRequestReq %v", err)
		render.ErrorJSON(c, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ct.UserFriendUsecase.RejectRequest(c, *req)
	if err != nil {
		logging.Errorf(c, "RejectRequest %v", err)
		switch err {
		case entity.ErrFriendRequestNotFound:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeUserFriendNotFound)
		case entity.ErrSentRequestNotFound:
			render.ErrorCodeJSON(c, err.Error(), http.StatusBadRequest, entity.CodeSentRequestNotFound)
		default:
			render.ErrorJSON(c, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	render.JSON(c, res)
}
