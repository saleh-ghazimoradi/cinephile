package gateway

import (
	"fmt"
	"github.com/saleh-ghazimoradi/cinephile/config"
	"github.com/saleh-ghazimoradi/cinephile/logger"
	"net/http"
	"os"
	"time"
)

func Server(router http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Appconfig.ServerAddress),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Logger.Info("starting server", "addr", srv.Addr, "env", config.Appconfig.Env)

	if err := srv.ListenAndServe(); err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}

	return nil
}
