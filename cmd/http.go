/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/saleh-ghazimoradi/cinephile/config"
	"github.com/saleh-ghazimoradi/cinephile/internal/gateway"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service"
	"github.com/saleh-ghazimoradi/cinephile/utils"
	"github.com/spf13/cobra"
	"log"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "launching the http rest listen server",
	Run: func(cmd *cobra.Command, args []string) {

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

		db, err := utils.PostgresConnection(cfg)
		if err != nil {
			log.Fatal(err)
		}

		movieDB := repository.NewMovieRepository(db)
		userDB := repository.NewUserRepository(db)

		movieService := service.NewMovieService(movieDB)
		userService := service.NewUserService(userDB)
		mailerService := service.NewMail(config.Appconfig.SMTPHost, config.Appconfig.SMTPPort, config.Appconfig.SMTPUsername, config.Appconfig.SMTPPassword, config.Appconfig.SMTPSender)

		movieHandler := gateway.NewMovieHandler(movieService)
		userHandler := gateway.NewUserHandler(userService, mailerService)

		routesHandler := gateway.Handlers{
			HealthCheckHandler:  gateway.HealthCheckHandler,
			ShowMovieHandler:    movieHandler.ShowMovieHandler,
			CreateMovieHandler:  movieHandler.CreateMovieHandler,
			UpdateMovieHandler:  movieHandler.UpdateMovieHandler,
			DeleteMovieHandler:  movieHandler.DeleteMovieHandler,
			ListMoviesHandler:   movieHandler.ListMoviesHandler,
			RegisterUserHandler: userHandler.RegisterUserHandler,
		}

		if err := gateway.Server(gateway.Routes(routesHandler)); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
