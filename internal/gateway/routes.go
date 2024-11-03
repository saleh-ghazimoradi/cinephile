package gateway

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handlers struct {
	HealthCheckHandler http.HandlerFunc
	ShowMovieHandler   http.HandlerFunc
	CreateMovieHandler http.HandlerFunc
}

func Routes(handlers Handlers) http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handlers.HealthCheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", handlers.ShowMovieHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", handlers.CreateMovieHandler)

	return router
}
