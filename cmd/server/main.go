package main

import "giiku-camp/internal/app"

func main() {
	app, err := app.New()
	if err != nil {
		panic(err)
	}
	defer app.Close()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
