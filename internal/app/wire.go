//go:build wireinject

package app

import (
	"giiku-camp/internal/app/config"
	"giiku-camp/internal/app/container"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func New() (*container.App, error) {
	wire.Build(
		provideGinEngine,
		config.New,

		container.NewApp,
	)
	return nil, nil
}

func provideGinEngine() *gin.Engine {
	return gin.Default()
}
