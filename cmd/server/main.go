package main

import (
	_ "giiku-camp/docs" // Swagger docs の読み込み
	"giiku-camp/internal/app"
)

// @title Giiku Camp API
// @version 1.0
// @description This is the Giiku Camp API server.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @query.collection.format multi
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	defer a.Close()

	if err := a.Run(); err != nil {
		panic(err)
	}
}
