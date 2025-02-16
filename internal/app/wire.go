//go:build wireinject

package app

import (
	"giiku-camp/internal/app/config"
	"giiku-camp/internal/app/container"
	"giiku-camp/internal/controller"
	repository "giiku-camp/internal/infra/datastore"
	"giiku-camp/internal/middleware"
	"giiku-camp/internal/usecase"
	"giiku-camp/pkg/datastore"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func New() (*container.App, error) {
	wire.Build(
		provideGinEngine,
		config.New,

		repository.NewUserRepo,
		repository.NewUserCrossingRepo,
		repository.NewUserFriendRepo,

		usecase.NewUserUsecase,
		usecase.NewUserCrossingUsecase,
		usecase.NewUserFriendUsecase,

		controller.NewUserCtrl,
		controller.NewUserCrossingCtrl,
		controller.NewUserFriendCtrl,

		middleware.NewMiddleware,

		datastore.New,
		container.NewCtrl,
		container.NewApp,
	)
	return nil, nil
}

func provideGinEngine() *gin.Engine {
	return gin.Default()
}
