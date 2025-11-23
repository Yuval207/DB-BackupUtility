package main

import (
	"fmt"
	"os"

	"github.com/antigravity/dbbackup/internal/config"
	"github.com/antigravity/dbbackup/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var appConfig config.Config

var rootCmd = &cobra.Command{
	Use:   "dbbackup",
	Short: "A CLI tool for database backups",
	Long:  `A comprehensive CLI tool for backing up and restoring databases (MySQL, PostgreSQL, MongoDB) to local or cloud storage.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dbbackup.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".dbbackup")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
		if err := viper.Unmarshal(&appConfig); err != nil {
			fmt.Printf("Unable to decode into struct, %v", err)
			os.Exit(1)
		}
		
		// Initialize logger
		logger.Init(appConfig.Log.Level)
	}
}

func main() {
	Execute()
}
