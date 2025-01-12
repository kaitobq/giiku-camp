package container

import (
	"fmt"
	"giiku-camp/internal/app/config"
	"giiku-camp/internal/controller"
	"giiku-camp/internal/infra/logging"
	"giiku-camp/internal/middleware"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

type container struct {
	userCtrl         controller.UserCtrl
	userCrossingCtrl controller.UserCrossingCtrl
}

func NewCtrl(userCtrl controller.UserCtrl, userCrossingCtrl controller.UserCrossingCtrl) container {
	return container{
		userCtrl:         userCtrl,
		userCrossingCtrl: userCrossingCtrl,
	}
}

type App struct {
	r          *gin.Engine
	cfg        *config.Config
	db         *datastore.Client
	middleware *middleware.Middleware
}

func NewApp(r *gin.Engine, container container, cfg *config.Config, db *datastore.Client, middleware *middleware.Middleware) *App {
	logging.Init()

	controller.SetUpRoutes(r, container.userCtrl, container.userCrossingCtrl, middleware)

	return &App{
		r:          r,
		cfg:        cfg,
		db:         db,
		middleware: middleware,
	}
}

func (a *App) Run() error {
	return a.r.Run(fmt.Sprintf(":%d", a.cfg.Port))
}

func (a *App) Close() {
	a.db.Close()
}
