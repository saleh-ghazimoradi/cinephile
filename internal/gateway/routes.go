package gateway

import "net/http"

type Handlers struct {
	HealthCheckHandler http.HandlerFunc
}

func Routes(handlers Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", handlers.HealthCheckHandler)

	return mux
}
