package eth

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

var client *ethclient.Client

func initClient() {
	var err error
	api_key := viper.GetString("infura.api_key")
	if api_key == "" {
		log.Fatal("/config/config.toml infura api key is required")
	}
	client, err = ethclient.Dial(fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", api_key))
	if err != nil {
		log.Fatal(err)
	}
}

func GetClient() *ethclient.Client {
	if client == nil {
		initClient()
	}
	return client
}
