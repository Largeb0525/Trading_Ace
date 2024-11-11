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
	client, err = ethclient.Dial(fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", viper.GetString("infura.api_key")))
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
