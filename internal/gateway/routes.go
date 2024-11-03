package gateway

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handlers struct {
	HealthCheckHandler http.HandlerFunc
}

func Routes(handlers Handlers) http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handlers.HealthCheckHandler)

	return router
}
