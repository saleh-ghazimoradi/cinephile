/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/saleh-ghazimoradi/cinephile/config"
	"github.com/saleh-ghazimoradi/cinephile/internal/gateway"
	"github.com/saleh-ghazimoradi/cinephile/logger"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "launching the http rest listen server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("server has started", "addr", config.Appconfig.ServerAddress, "env", config.Appconfig.Env)

		movieHandler := gateway.NewMovieHandler()

		routesHandler := gateway.Handlers{
			HealthCheckHandler: gateway.HealthCheckHandler,
			ShowMovieHandler:   movieHandler.ShowMovieHandler,
			CreateMovieHandler: movieHandler.CreateMovieHandler,
		}

		if err := gateway.Server(gateway.Routes(routesHandler)); err != nil {
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
