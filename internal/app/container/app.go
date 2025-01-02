package container

import (
	"fmt"
	"giiku-camp/internal/app/config"
	"giiku-camp/internal/controller"
	"giiku-camp/internal/infra/logging"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

type container struct {
	userCtrl controller.UserCtrl
}

func NewCtrl(userCtrl controller.UserCtrl) container {
	return container{
		userCtrl: userCtrl,
	}
}

type App struct {
	r   *gin.Engine
	cfg *config.Config
	db  *firestore.Client
}

func NewApp(r *gin.Engine, container container, cfg *config.Config, db *firestore.Client) *App {
	logging.Init()

	controller.SetUpRoutes(r, container.userCtrl)

	return &App{
		r:   r,
		cfg: cfg,
		db:  db,
	}
}

func (a *App) Run() error {
	return a.r.Run(fmt.Sprintf(":%d", a.cfg.Port))
}

func (a *App) Close() {
	a.db.Close()
}
