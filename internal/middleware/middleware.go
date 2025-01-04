package middleware

import "giiku-camp/internal/domain/repository"

type Middleware struct {
	API *API
}

func NewMiddleware(userRepo repository.UserRepo) *Middleware {
	return &Middleware{
		API: NewAPI(userRepo),
	}
}
