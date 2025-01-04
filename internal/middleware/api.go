package middleware

import (
	"giiku-camp/internal/domain/repository"
	"giiku-camp/internal/domain/xcontext"
	"giiku-camp/internal/infra/logging"
	"giiku-camp/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type API struct {
	userRepo repository.UserRepo
}

func NewAPI(userRepo repository.UserRepo) *API {
	return &API{
		userRepo: userRepo,
	}
}

func (a *API) withUser(c *gin.Context) error {
	isValid, err := jwt.VerifyReqHeaderToken(c)
	if err != nil || !isValid {
		logging.Infof(c, "middleware.API.withUser jwt.VerifyReqHeaderToken %v", err)
		// render.ErrorJSON(c, err.Error(), http.StatusUnauthorized)
		return err
	}

	userID, err := jwt.ExtractUserIDFromToken(c)
	if err != nil {
		logging.Infof(c, "middleware.API.withUser jwt.ExtractUserIDFromToken %v", err)
		// render.ErrorJSON(c, err.Error(), http.StatusUnauthorized)
		return err
	}

	user, err := a.userRepo.FindByID(c, userID)
	if err != nil {
		logging.Infof(c, "middleware.API.withUser a.userRepo.FindByID %v", err)
		// render.ErrorJSON(c, err.Error(), http.StatusUnauthorized)
		return err
	}

	xcontext.WithUser(c, user)
	return nil
}

func (a *API) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := a.withUser(c); err != nil {
			// render.ErrorJSON(c, err.Error(), http.StatusUnauthorized)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
