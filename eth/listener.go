package eth

import (
	"context"
	"fmt"
	"log"

	"github.com/Largeb0525/Trading_Ace/database"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

func ListenToContractEvents() {
	client := GetClient()
	campaignAddresses, err := database.GetActiveCampaignAddresses()
	if err != nil {
		log.Fatal(err)
	}
	if len(campaignAddresses) > 0 {
		AddAddresses(campaignAddresses)
	}
	for {
		addrList := GetAddresses()
		if len(addrList) == 0 {
			fmt.Println("No addresses to listen to, waiting for changes...")
			<-GetNotifyChannel()
			continue
		}

		query := ethereum.FilterQuery{Addresses: addrList}
		logs := make(chan types.Log)
		ctx, cancel := context.WithCancel(context.Background())
		sub, err := client.SubscribeFilterLogs(ctx, query, logs)
		if err != nil {
			log.Printf("Failed to subscribe: %v", err)
			cancel()
			<-GetNotifyChannel()
			continue
		}
		fmt.Println("Listening for Swap events...")

		done := make(chan struct{})
		go func() {
			for {
				select {
				case err := <-sub.Err():
					if err != nil {
						log.Printf("Subscription error: %v", err)
					}
					cancel()
					done <- struct{}{}
					return
				case vLog := <-logs:
					handleLogs(vLog)
				}
			}
		}()

		<-GetNotifyChannel()
		cancel()
		sub.Unsubscribe()
		<-done
	}
}
