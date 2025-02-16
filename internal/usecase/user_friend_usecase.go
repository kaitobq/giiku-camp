package usecase

import (
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"
	"giiku-camp/internal/domain/xcontext"
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

type userFriendUsecase struct {
	userFriendRepo repository.UserFriendRepo
	userRepo       repository.UserRepo
	db             *datastore.Client
}

func NewUserFriendUsecase(userFriendRepo repository.UserFriendRepo, userRepo repository.UserRepo, db *datastore.Client) UserFriendUsecase {
	return &userFriendUsecase{userFriendRepo: userFriendRepo, userRepo: userRepo, db: db}
}

func (u *userFriendUsecase) GetUserFriend(c *gin.Context) (*response.UserFriendRes, error) {
	user := xcontext.User(c)
	userFriends, err := u.userFriendRepo.FindByUserID(c.Request.Context(), user.ID)
	if err != nil {
		return nil, err
	}

	var users []entity.User
	for _, userFriend := range userFriends {
		user, err := u.userRepo.FindByID(c.Request.Context(), userFriend.RelatedUserID)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return response.NewUserFriendRes(userFriends, users)
}

func (u *userFriendUsecase) SendRequest(c *gin.Context, req request.SendRequestReq) (*response.SendRequestRes, error) {
	sender := xcontext.User(c)
	reciever, err := u.userRepo.FindByID(c.Request.Context(), req.UserID)
	if err != nil {
		return nil, err
	}

	senderEntList, err := u.userFriendRepo.FindByUserID(c.Request.Context(), sender.ID)
	if err != nil {
		switch err {
		case entity.ErrUserFriendNotFound:
			senderEntList = []entity.UserFriend{}
		default:
			return nil, err
		}
	}
	recieverEntList, err := u.userFriendRepo.FindByUserID(c.Request.Context(), reciever.ID)
	if err != nil {
		switch err {
		case entity.ErrUserFriendNotFound:
			recieverEntList = []entity.UserFriend{}
		default:
			return nil, err
		}
	}

	for _, senderEnt := range senderEntList {
		if senderEnt.IsFriend(reciever.ID) {
			return nil, entity.ErrAlreadyFriend
		} else if senderEnt.IsSentRequest(reciever.ID) {
			return nil, entity.ErrFriendRequestAlreadySent
		} else if senderEnt.IsReceivedRequest(reciever.ID) {
			return nil, entity.ErrFriendRequestAlreadyReceived
		}
	}
	for _, recieverEnt := range recieverEntList {
		if recieverEnt.IsFriend(sender.ID) {
			return nil, entity.ErrAlreadyFriend
		} else if recieverEnt.IsSentRequest(sender.ID) {
			return nil, entity.ErrFriendRequestAlreadySent
		} else if recieverEnt.IsReceivedRequest(sender.ID) {
			return nil, entity.ErrFriendRequestAlreadyReceived
		}
	}

	senderEnt := entity.NewUserFriend(sender.ID, reciever.ID, entity.UserFriendTypeSentRequest)
	recieverEnt := entity.NewUserFriend(reciever.ID, sender.ID, entity.UserFriendTypeReceivedRequest)

	_, err = u.db.RunInTransaction(c.Request.Context(), func(tx *datastore.Transaction) error {
		if err := u.userFriendRepo.UpdateWithTransaction(tx, senderEnt); err != nil {
			return err
		}
		if err := u.userFriendRepo.UpdateWithTransaction(tx, recieverEnt); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response.NewSendRequestRes()
}

func (u *userFriendUsecase) AcceptRequest(c *gin.Context, req request.AcceptRequestReq) (*response.AcceptRequestRes, error) {
	accepter := xcontext.User(c)
	sender, err := u.userRepo.FindByID(c.Request.Context(), req.UserID)
	if err != nil {
		return nil, err
	}

	accepterEnt, err := u.userFriendRepo.FindByRelatedUserID(c.Request.Context(), accepter.ID)
	if err != nil {
		return nil, err
	}

	senderEnt, err := u.userFriendRepo.FindByRelatedUserID(c.Request.Context(), sender.ID)
	if err != nil {
		return nil, err
	}

	if !accepterEnt.IsReceivedRequest(sender.ID) {
		return nil, entity.ErrFriendRequestNotFound
	}

	if !senderEnt.IsSentRequest(accepter.ID) {
		return nil, entity.ErrSentRequestNotFound
	}

	accepterEnt.AcceptRequest()
	senderEnt.AcceptedRequest()

	_, err = u.db.RunInTransaction(c.Request.Context(), func(tx *datastore.Transaction) error {
		if err := u.userFriendRepo.UpdateWithTransaction(tx, accepterEnt); err != nil {
			return err
		}
		if err := u.userFriendRepo.UpdateWithTransaction(tx, senderEnt); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response.NewAcceptRequestRes()
}

func (u *userFriendUsecase) RejectRequest(c *gin.Context, req request.RejectRequestReq) (*response.RejectRequestRes, error) {
	rejecter := xcontext.User(c)
	sender, err := u.userRepo.FindByID(c.Request.Context(), req.UserID)
	if err != nil {
		return nil, err
	}

	rejecterEnt, err := u.userFriendRepo.FindByRelatedUserID(c.Request.Context(), rejecter.ID)
	if err != nil {
		return nil, err
	}
	senderEnt, err := u.userFriendRepo.FindByRelatedUserID(c.Request.Context(), sender.ID)
	if err != nil {
		return nil, err
	}

	if !rejecterEnt.IsReceivedRequest(sender.ID) {
		return nil, entity.ErrFriendRequestNotFound
	}
	if !senderEnt.IsSentRequest(rejecter.ID) {
		return nil, entity.ErrSentRequestNotFound
	}

	_, err = u.db.RunInTransaction(c.Request.Context(), func(tx *datastore.Transaction) error {
		if err := u.userFriendRepo.DeleteWithTransaction(tx, rejecterEnt); err != nil {
			return err
		}
		if err := u.userFriendRepo.DeleteWithTransaction(tx, senderEnt); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response.NewRejectRequestRes()
}
