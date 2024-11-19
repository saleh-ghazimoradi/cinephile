package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/cinephile/config"
	"github.com/saleh-ghazimoradi/cinephile/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		logger.Logger.Info("completing background tasks", "addr", srv.Addr)

		wg.Wait()
		shutdownError <- nil
	}()

	logger.Logger.Info("starting server", "addr", srv.Addr, "env", config.Appconfig.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	logger.Logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
