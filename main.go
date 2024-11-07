package main

import (
	"fmt"
	"os"

	"github.com/Largeb0525/Trading_Ace/cmd"
	"github.com/Largeb0525/Trading_Ace/internal"
	"github.com/spf13/viper"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	router := internal.InitRouter()
	err = router.Run(":" + port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
