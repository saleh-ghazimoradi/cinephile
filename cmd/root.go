/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/saleh-ghazimoradi/cinephile/config"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cinephile",
	Short: "A simple movie recommendation app",
	Long: `Cinephile is a simple Golang backend with PostgreSQL, featuring JWT authentication, 
role-based access, and a clean RESTful API for efficient data management and seamless integration.`,
}

func Execute() {
	err := os.Setenv("TZ", time.UTC.String())
	if err != nil {
		panic(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		initConfig,
	)
}

func initConfig() {
	err := config.LoadingConfig(".")
	if err != nil {
		log.Fatal(err)
	}
}
