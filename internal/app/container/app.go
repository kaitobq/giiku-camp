package container

import (
	"fmt"
	"giiku-camp/internal/app/config"
	"giiku-camp/internal/ctrl"

	"github.com/gin-gonic/gin"
)

type container struct {
}

func NewCtrl() *container {
	return &container{}
}

type App struct {
	r   *gin.Engine
	cfg *config.Config
}

func NewApp(r *gin.Engine, cfg *config.Config) *App {
	ctrl.SetUpRoutes(r)

	return &App{
		r:   r,
		cfg: cfg,
	}
}

func (a *App) Run() error {
	return a.r.Run(fmt.Sprintf(":%d", a.cfg.Port))
}
