//go:build wireinject

package app

import (
	"giiku-camp/internal/app/config"
	"giiku-camp/internal/app/container"
	"giiku-camp/internal/controller"
	repository "giiku-camp/internal/infra/firestore"
	"giiku-camp/internal/middleware"
	"giiku-camp/internal/usecase"
	"giiku-camp/pkg/firestore"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func New() (*container.App, error) {
	wire.Build(
		provideGinEngine,
		config.New,

		repository.NewUserRepo,

		usecase.NewUserUsecase,

		controller.NewUserCtrl,

		middleware.NewMiddleware,

		firestore.New,
		container.NewCtrl,
		container.NewApp,
	)
	return nil, nil
}

func provideGinEngine() *gin.Engine {
	return gin.Default()
}
