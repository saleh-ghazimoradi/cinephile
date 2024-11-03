package gateway

import (
	"github.com/saleh-ghazimoradi/cinephile/config"
	"net/http"
)

const version = "1.0.0"

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": config.Appconfig.Env,
			"version":     version,
		},
	}

	if err := writeJSON(w, http.StatusOK, env, nil); err != nil {
		serverErrorResponse(w, r, err)
	}
}
