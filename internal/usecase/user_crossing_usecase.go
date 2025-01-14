package usecase

import (
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"
	"giiku-camp/internal/domain/xcontext"
	"giiku-camp/internal/infra/logging"
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"

	"github.com/gin-gonic/gin"
)

type userCrossingUsecase struct {
	userCrossingRepo repository.UserCrossingRepo
	userRepo         repository.UserRepo
}

func NewUserCrossingUsecase(userCrossingRepo repository.UserCrossingRepo, userRepo repository.UserRepo) UserCrossingUsecase {
	return &userCrossingUsecase{userCrossingRepo: userCrossingRepo, userRepo: userRepo}
}

func (u *userCrossingUsecase) RegisterUserCrossing(c *gin.Context, req request.RegisterUserCrossingReq) (*response.RegisterUserCrossingRes, error) {
	user := xcontext.User(c)
	var users []entity.User
	for _, id := range req.UserIDs {
		userCrossing := entity.NewUserCrossing(user.ID, id)

		userCrossing.UpdateCreatedAt()
		ctx := c.Request.Context()
		err := u.userCrossingRepo.Update(ctx, *userCrossing)
		if err != nil {
			logging.Errorf(c, "Update %v", err)
			return nil, err
		}

		us, err := u.userRepo.FindByID(ctx, id)
		if err != nil {
			logging.Errorf(c, "FindByID %v", err)
			return nil, err
		}

		users = append(users, *us)
	}

	return response.NewRegisterUserCrossingRes(users)
}

func (u *userCrossingUsecase) GetUserCrossing(c *gin.Context) (*response.GetUserCrossingRes, error) {
	user := xcontext.User(c)
	logging.Infof(c, "user: %v", user)
	for _, id := range user.UserCrossingIDs {
		ctx := c.Request.Context()
		us, err := u.userRepo.FindByID(ctx, id)
		if err != nil {
			logging.Errorf(c, "FindByID %v", err)
			return nil, err
		}
		if user == us {
			continue
		}
		res, err := response.NewGetUserCrossingRes([]entity.User{*user})
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
