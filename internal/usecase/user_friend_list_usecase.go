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

type userFriendListUsecase struct {
	userFriendListRepo repository.UserFriendListRepo
	userRepo           repository.UserRepo
	db                 *datastore.Client
}

func NewUserFriendListUsecase(userFriendListRepo repository.UserFriendListRepo, userRepo repository.UserRepo, db *datastore.Client) UserFriendListUsecase {
	return &userFriendListUsecase{userFriendListRepo: userFriendListRepo, userRepo: userRepo, db: db}
}

func (u *userFriendListUsecase) GetUserFriendList(c *gin.Context) (*response.UserFriendListRes, error) {
	user := xcontext.User(c)
	userFriendList, err := u.userFriendListRepo.FindByUserID(c.Request.Context(), user.ID)
	if err != nil {
		if err == entity.ErrUserFriendListNotFound { // TODO: 全データ削除するタイミングで消す
			ent := entity.NewUserFriendList(user.ID)
			if err := u.userFriendListRepo.Update(c.Request.Context(), ent); err != nil {
				return nil, err
			}
			return response.NewUserFriendListRes(*ent), nil
		}
		return nil, err
	}
	return response.NewUserFriendListRes(*userFriendList), nil
}

func (u *userFriendListUsecase) SendRequest(c *gin.Context, req request.SendRequestReq) (*response.SendRequestRes, error) {
	sender := xcontext.User(c)
	ctx := c.Request.Context()
	receiver, err := u.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	senderEnt, err := u.userFriendListRepo.FindByUserID(ctx, sender.ID)
	if err != nil {
		if err == entity.ErrUserFriendListNotFound { // TODO: 全データ削除するタイミングで消す
			senderEnt = entity.NewUserFriendList(sender.ID)
			if err := u.userFriendListRepo.Update(ctx, senderEnt); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	receiverEnt, err := u.userFriendListRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		if err == entity.ErrUserFriendListNotFound { // TODO: 全データ削除するタイミングで消す
			receiverEnt = entity.NewUserFriendList(req.UserID)
			if err := u.userFriendListRepo.Update(ctx, receiverEnt); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// senderが既にリクエストを送っている
	if senderEnt.HasSentRequest(receiver) || receiverEnt.HasFriendRequest(sender) {
		return nil, entity.ErrFriendRequestAlreadySent
	}

	// receiverから既にリクエストを受けている
	if senderEnt.HasFriendRequest(receiver) || receiverEnt.HasSentRequest(sender) {
		return nil, entity.ErrFriendRequestAlreadyReceived
	}

	// 既にフレンド
	if senderEnt.HasFriend(receiver) || receiverEnt.HasFriend(sender) {
		return nil, entity.ErrAlreadyFriend
	}

	senderEnt.AddSentRequest(receiver)
	receiverEnt.AddFriendRequest(sender)

	// 送信者、受信者の整合性を担保する
	_, err = u.db.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := u.userFriendListRepo.UpdateWithTransaction(tx, senderEnt); err != nil {
			return err
		}
		if err := u.userFriendListRepo.UpdateWithTransaction(tx, receiverEnt); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response.NewSendRequestRes()
}

func (u *userFriendListUsecase) AcceptRequest(c *gin.Context, req request.AcceptRequestReq) (*response.AcceptRequestRes, error) {
	accepter := xcontext.User(c)
	ctx := c.Request.Context()
	sender, err := u.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	accepterEnt, err := u.userFriendListRepo.FindByUserID(ctx, accepter.ID)
	if err != nil {
		return nil, err
	}
	senderEnt, err := u.userFriendListRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// if accepterEnt.HasFriend(req.UserID) {
	// 	return nil, entity.ErrAlreadyFriend
	// }
	if !accepterEnt.HasFriendRequest(sender) {
		return nil, entity.ErrFriendRequestNotFound
	}
	if !senderEnt.HasSentRequest(accepter) {
		return nil, entity.ErrSentRequestNotFound
	}

	accepterEnt.AddFriend(sender)
	accepterEnt.RemoveFriendRequest(sender)
	senderEnt.AddFriend(accepter)
	senderEnt.RemoveSentRequest(accepter)

	// 送信者、受信者の整合性を担保する
	_, err = u.db.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := u.userFriendListRepo.UpdateWithTransaction(tx, accepterEnt); err != nil {
			return err
		}
		if err := u.userFriendListRepo.UpdateWithTransaction(tx, senderEnt); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response.NewAcceptRequestRes()
}

func (u *userFriendListUsecase) RejectRequest(c *gin.Context, req request.RejectRequestReq) (*response.RejectRequestRes, error) {
	rejecter := xcontext.User(c)
	ctx := c.Request.Context()
	sender, err := u.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	rejecterEnt, err := u.userFriendListRepo.FindByUserID(ctx, rejecter.ID)
	if err != nil {
		return nil, err
	}
	senderEnt, err := u.userFriendListRepo.FindByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// if rejecterEnt.HasFriend(req.UserID) {
	// 	return nil, entity.ErrAlreadyFriend
	// }
	if !rejecterEnt.HasFriendRequest(sender) {
		return nil, entity.ErrFriendRequestNotFound
	}
	if !senderEnt.HasSentRequest(rejecter) {
		return nil, entity.ErrSentRequestNotFound
	}

	rejecterEnt.RemoveFriendRequest(sender)
	senderEnt.RemoveSentRequest(rejecter)

	// 送信者、受信者の整合性を担保する
	_, err = u.db.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := u.userFriendListRepo.UpdateWithTransaction(tx, rejecterEnt); err != nil {
			return err
		}
		if err := u.userFriendListRepo.UpdateWithTransaction(tx, senderEnt); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response.NewRejectRequestRes()
}
