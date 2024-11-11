package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Largeb0525/Trading_Ace/database"
	"github.com/Largeb0525/Trading_Ace/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "trading_ace",
	Short: "trading_ace",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.InitPostgreSQL()
		defer db.Close()
		server.StartServer()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("config", "", "config file (default is ./config.toml)")
}

func initConfig() {
	configFile, _ := rootCmd.Flags().GetString("config")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")
		viper.SetConfigType("toml")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
			os.Exit(1)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %s", err)
	}
}
