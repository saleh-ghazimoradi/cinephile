/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/cinephile/config"
	"github.com/saleh-ghazimoradi/cinephile/internal/gateway"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service"
	"github.com/saleh-ghazimoradi/cinephile/logger"
	"github.com/saleh-ghazimoradi/cinephile/utils"
	"github.com/spf13/cobra"
	"log"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "launching the http rest listen server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("server has started", "addr", config.Appconfig.ServerAddress, "env", config.Appconfig.Env)

		cfg := utils.PostgresConfig{
			Host:         config.Appconfig.DBHost,
			Port:         config.Appconfig.DBPort,
			User:         config.Appconfig.DBUser,
			Password:     config.Appconfig.DBPassword,
			Database:     config.Appconfig.DBName,
			SSLMode:      config.Appconfig.DBSSLMode,
			MaxIdleTime:  config.Appconfig.MaxIdleTime,
			MaxIdleConns: config.Appconfig.MaxIdleConns,
			MaxOpenConns: config.Appconfig.MaxOpenConns,
			Timeout:      config.Appconfig.Timeout,
		}
		fmt.Println(cfg)

		db, err := utils.PostgresConnection(cfg)
		if err != nil {
			logger.Logger.Error(err.Error())
		}

		movieDB := repository.NewMovieRepository(db)
		movieService := service.NewMovieService(movieDB)
		movieHandler := gateway.NewMovieHandler(movieService)

		routesHandler := gateway.Handlers{
			HealthCheckHandler: gateway.HealthCheckHandler,
			ShowMovieHandler:   movieHandler.ShowMovieHandler,
			CreateMovieHandler: movieHandler.CreateMovieHandler,
		}

		if err := gateway.Server(gateway.Routes(routesHandler)); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
