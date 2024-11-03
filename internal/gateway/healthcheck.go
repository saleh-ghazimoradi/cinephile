package gateway

import (
	"fmt"
	"github.com/saleh-ghazimoradi/cinephile/config"
	"net/http"
)

const version = "1.0.0"

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "enviroment:%s\n", config.Appconfig.Env)
	fmt.Fprintln(w, "version: "+version)
}
