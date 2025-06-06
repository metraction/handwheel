/*
Copyright 2025 NAME HERE EMAIL ADDRESS
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/metraction/handwheel/model"
	"github.com/metraction/handwheel/routing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config *model.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ai-backend-api",
	Short: "Backend for mentor mate",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Start health and readiness probe server
		healthServerPort := config.HttpServer.Port
		if healthServerPort == "" {
			healthServerPort = "8081"
			log.Printf("Health server port not set in config, defaulting to %s", healthServerPort)
		}
		go routing.StartHealthServer(healthServerPort)

		// Start Prometheus fetcher in the background
		routing.ProtheusCraneDevLakeRouter(config)

		fmt.Println("Started handler and http server.")
		select {}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".image-handler" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".") // Add current directory to config search path
		viper.SetConfigName(".image-handler")
	}
	viper.AutomaticEnv() // read in environment variables that match

	// Open config file for ENV variables substitution
	file, err := os.Open(viper.ConfigFileUsed())
	if err != nil {
		log.Fatal("No config file found ", err)
		config = &model.Config{}
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading config file", err)
	}
	expandedContent := os.ExpandEnv(string(content))
	myReader := strings.NewReader(expandedContent)
	// If a config file is found, read it in.
	if err := viper.ReadConfig(myReader); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Error loading config", err)
		config = &model.Config{}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Unable to decode config into struct", err)
	}

	// Create logic

}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mentor-mate.yaml)")
}
