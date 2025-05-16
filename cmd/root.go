/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Tiktai/handler/model"

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
		fmt.Println("Starting HTTP server on :8080 ...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
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

		// Search config in home directory with name ".mx" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mentor-mate")
	}
	viper.AutomaticEnv() // read in environment variables that match

	// Open config file for ENV variables substitution
	file, err := os.Open(viper.ConfigFileUsed())
	if err != nil {
		log.Fatal("No config file found", err)
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
		fmt.Println("No config file found")
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
