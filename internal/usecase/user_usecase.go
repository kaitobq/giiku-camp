package usecase

import (
	"giiku-camp/internal/domain/entity"
	"giiku-camp/internal/domain/repository"
	"giiku-camp/internal/domain/xcontext"
	"giiku-camp/internal/infra/apple"
	"giiku-camp/internal/infra/logging"
	"giiku-camp/internal/usecase/request"
	"giiku-camp/internal/usecase/response"
	"giiku-camp/pkg/jwt"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

type userUsecase struct {
	userRepo           repository.UserRepo
	userFriendListRepo repository.UserFriendListRepo
	appleClient        apple.Client
	db                 *datastore.Client
}

func NewUserUsecase(userRepo repository.UserRepo, userFriendListRepo repository.UserFriendListRepo, appleClient apple.Client, db *datastore.Client) UserUsecase {
	return &userUsecase{userRepo: userRepo, userFriendListRepo: userFriendListRepo, appleClient: appleClient, db: db}
}

func (u *userUsecase) SignUp(c *gin.Context, req request.SignUpReq) (*response.SignUpRes, error) {
	user, err := entity.NewUser(req.Name, &req.GitHubID, &req.QiitaID, &req.ZennID, &req.XID)
	if err != nil {
		logging.Errorf(c, "NewUser %v", err)
		return nil, err
	}

	ctx := c.Request.Context()
	res, err := u.appleClient.VerifyAuthorizationCode(ctx, req.AuthorizationCode)
	if err != nil {
		logging.Errorf(c, "VerifyAuthorizationCode %v", err)
		return nil, err
	}
	exist, err := u.userRepo.FindByAppleID(ctx, res.UserID)
	switch err {
	case nil:
		if exist != nil {
			logging.Infof(c, "FindByAppleID returned ErrAppleIDAlreadyUsed(user: %v)", exist)
			return nil, entity.ErrAppleIDAlreadyUsed
		}
	case entity.ErrUserNotFound:
		user.SetAppleID(res.UserID)
	default:
		logging.Errorf(c, "FindByAppleID %v", err)
		return nil, err
	}

	_, err = u.db.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		user.UpdateCreatedAt()
		if err := u.userRepo.UpdateWithTransaction(tx, user); err != nil {
			logging.Errorf(c, "Update %v", err)
			return err
		}
		friendList := entity.NewUserFriendList(user.ID)
		if err := u.userFriendListRepo.UpdateWithTransaction(tx, friendList); err != nil {
			logging.Errorf(c, "Update %v", err)
			return err
		}

		return nil
	})
	if err != nil {
		logging.Errorf(c, "RunInTransaction %v", err)
		return nil, err
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID)
	if err != nil {
		logging.Errorf(c, "GenerateAccessToken %v", err)
		return nil, err
	}
	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.TokenVersion)
	if err != nil {
		logging.Errorf(c, "GenerateRefreshToken %v", err)
		return nil, err
	}

	return response.NewSignUpRes(user, accessToken, refreshToken)
}

func (u *userUsecase) SignIn(c *gin.Context, req request.SignInReq) (*response.SignInRes, error) {
	ctx := c.Request.Context()

	res, err := u.appleClient.VerifyAuthorizationCode(ctx, req.AuthorizationCode)
	if err != nil {
		logging.Errorf(c, "VerifyAuthorizationCode %v", err)
		return nil, err
	}
	user, err := u.userRepo.FindByAppleID(ctx, res.UserID)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			logging.Infof(c, "FindByAppleID returned ErrUserNotFound(appleID: %s)", res.UserID)
			return nil, entity.ErrUserNotFound
		default:
			logging.Errorf(c, "FindByAppleID %v", err)
			return nil, err
		}
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID)
	if err != nil {
		logging.Errorf(c, "GenerateAccessToken %v", err)
		return nil, err
	}
	user.IncrementTokenVersion()
	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.TokenVersion)
	if err != nil {
		logging.Errorf(c, "GenerateRefreshToken %v", err)
		return nil, err
	}
	if err = u.userRepo.Update(ctx, user); err != nil {
		logging.Errorf(c, "Update %v", err)
		return nil, err
	}

	return response.NewSignInRes(user, accessToken, refreshToken)
}

func (u *userUsecase) RefreshToken(c *gin.Context, req request.RefreshTokenReq) (*response.RefreshTokenRes, error) {
	token, err := jwt.VerifyToken(req.RefreshToken)
	if err != nil {
		logging.Errorf(c, "VerifyToken %v", err)
		return nil, err
	}

	userID, err := jwt.ExtractUserIDFromToken(token)
	if err != nil {
		logging.Errorf(c, "ExtractUserIDFromToken %v", err)
		return nil, err
	}

	ctx := c.Request.Context()
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logging.Errorf(c, "FindByID %v", err)
		return nil, err
	}

	accessToken, refreshToken, err := jwt.RefreshTokens(*user, req.RefreshToken)
	if err != nil {
		logging.Errorf(c, "RefreshTokens %v", err)
		return nil, err
	}
	user.IncrementTokenVersion()
	if err = u.userRepo.Update(ctx, user); err != nil {
		logging.Errorf(c, "Update %v", err)
		return nil, err
	}

	return response.NewRefreshTokenRes(accessToken, refreshToken)
}

func (u *userUsecase) GetMe(c *gin.Context) (*response.UserRes, error) {
	user := xcontext.User(c)
	return response.NewUserRes(user), nil
}

func (u *userUsecase) GetUser(c *gin.Context, userID string) (*response.CrossedUserRes, error) {
	ctx := c.Request.Context()
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logging.Errorf(c, "FindByID %v", err)
		return nil, err
	}

	res := response.NewCrossedUserRes(*user)
	return &res, nil
}

func (u *userUsecase) UpdateMe(c *gin.Context, req request.UpdateMeReq) (*response.UserRes, error) {
	user := xcontext.User(c)
	updates := map[string]*string{
		"Name":        &req.Name,
		"Description": &req.Description,
		"GitHubID":    &req.GitHubID,
		"QiitaID":     &req.QiitaID,
		"ZennID":      &req.ZennID,
		"XID":         &req.XID,
	}

	for field, value := range updates {
		if isChanged(*value) {
			switch field {
			case "Name":
				user.Name = *value
			case "Description":
				user.Description = *value
			case "GitHubID":
				user.GitHubID = *value
			case "QiitaID":
				user.QiitaID = *value
			case "ZennID":
				user.ZennID = *value
			case "XID":
				user.XID = *value
			}
		}
	}

	ctx := c.Request.Context()
	if err := u.userRepo.Update(ctx, user); err != nil {
		logging.Errorf(c, "Update %v", err)
		return nil, err
	}

	return response.NewUserRes(user), nil
}

func isChanged(s string) bool {
	return s != ""
}
