package eth

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/Largeb0525/Trading_Ace/database"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	swapEventABI = `[{
		"anonymous": false,
		"inputs": [
		  {"indexed": true, "internalType": "address", "name": "sender", "type": "address"},
		  {"indexed": false, "internalType": "uint256", "name": "amount0In", "type": "uint256"},
		  {"indexed": false, "internalType": "uint256", "name": "amount1In", "type": "uint256"},
		  {"indexed": false, "internalType": "uint256", "name": "amount0Out", "type": "uint256"},
		  {"indexed": false, "internalType": "uint256", "name": "amount1Out", "type": "uint256"},
		  {"indexed": true, "internalType": "address", "name": "to", "type": "address"}
		],
		"name": "Swap",
		"type": "event"
	}]`
	swapEventTopicHash = "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"
)

type SwapEvent struct {
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
}

type SwapInfo struct {
	Sender      string
	USDC        float64
	Timestamp   int64
	PoolAddress string
	TxHash      string
}

func getBlockTime(blockNumber uint64) (int64, error) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return 0, err
	}
	return int64(block.Time()), nil
}

func timeToBlockNumber(targetTime int64, afterBlock uint64, beforeBlock uint64) (uint64, error) {
	var low, high uint64
	if afterBlock != 0 {
		low = afterBlock
	} else {
		low = 0
	}
	if beforeBlock != 0 {
		high = beforeBlock
	} else {
		latestBlock, err := client.BlockByNumber(context.Background(), nil)
		if err != nil {
			return 0, err
		}
		high = latestBlock.NumberU64()
	}

	for low <= high {
		mid := (low + high) / 2
		if mid == low {
			break
		}
		blockTime, err := getBlockTime(mid)
		if err != nil {
			return 0, err
		}
		if blockTime == targetTime {
			return mid, nil
		} else if blockTime < targetTime {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if low == 0 {
		return 0, nil
	}
	return low - 1, nil
}

func FetchSwapEvents(poolAddress string, startTime, endTime int64) ([]types.Log, error) {
	startBlock, err := timeToBlockNumber(startTime, 0, 0)
	if err != nil {
		return nil, err
	}
	endBlock, err := timeToBlockNumber(endTime, startBlock, 0)
	if err != nil {
		return nil, err
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startBlock + 1)),
		ToBlock:   big.NewInt(int64(endBlock)),
		Addresses: []common.Address{common.HexToAddress(poolAddress)},
		Topics:    [][]common.Hash{{common.HexToHash(swapEventTopicHash)}},
	}

	return client.FilterLogs(context.Background(), query)
}

func ParseSwapEvents(logs []types.Log) ([]SwapInfo, error) {
	var swapInfos []SwapInfo
	contractABI, err := abi.JSON(strings.NewReader(swapEventABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	for _, vLog := range logs {
		if len(vLog.Topics) != 3 ||
			vLog.Topics[0].Hex() != swapEventTopicHash {
			continue
		}

		var swapInfo SwapInfo
		swapInfo.Sender = ParseAddress(vLog.Topics[1].Hex())
		swapInfo.PoolAddress = vLog.Address.Hex()
		swapInfo.TxHash = vLog.TxHash.Hex()
		t, err := getBlockTime(vLog.BlockNumber)
		if err != nil {
			log.Printf("Failed to get log timestamp: %v", err)
			continue
		}
		swapInfo.Timestamp = t

		var swapEvent SwapEvent
		err = contractABI.UnpackIntoInterface(&swapEvent, "Swap", vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack Swap event: %v", err)
			continue
		}
		if swapEvent.Amount0In.Cmp(big.NewInt(0)) > 0 {
			swapInfo.USDC = float64(swapEvent.Amount0In.Int64()) / 1e6
		} else if swapEvent.Amount0Out.Cmp(big.NewInt(0)) > 0 {
			swapInfo.USDC = float64(swapEvent.Amount0Out.Int64()) / 1e6
		} else {
			log.Printf("Invalid USDC value in log: %v", vLog)
			continue
		}

		swapInfos = append(swapInfos, swapInfo)
	}

	return swapInfos, nil
}

func ProcessSwapInfos(task database.Task, swaps []SwapInfo, onboardingTask database.Task) {
	senderMap := make(map[string]float64)
	for _, swap := range swaps {
		senderMap[swap.Sender] += swap.USDC
	}

	validatedSenderMap, totalUSDC := calculateTotalUSDC(senderMap, onboardingTask.TaskID, onboardingTask.OnboardingThreshold)

	for sender, usdc := range validatedSenderMap {
		reward := task.PointsPool * (usdc / totalUSDC)
		reward = math.Floor(reward*1e6) / 1e6

		user, err := database.GetUserByAddress(sender)
		if err != nil {
			log.Printf("Failed to get user by address: %v", err)
			continue
		}

		err = database.UpdateUserTaskByUserIDTaskID(user.UserID, task.TaskID, true, usdc, reward)
		if err != nil {
			log.Printf("Failed to update user task: %v", err)
			continue
		}
		if reward > 0 {
			err = database.CreateUserPointsHistory(user.UserID, task.TaskID, task.CampaignID, reward)
			if err != nil {
				log.Printf("Failed to create user points history: %v", err)
				continue
			}
		}
	}
}

func calculateTotalUSDC(senderMap map[string]float64, taskID int, threshold float64) (map[string]float64, float64) {
	totalUSDC := 0.0

	for sender, usdc := range senderMap {
		if usdc >= threshold {
			totalUSDC += usdc
			continue
		}
		user, err := database.GetUserByAddress(sender)
		if err != nil {
			log.Printf("Failed to get user by address: %v", err)
			continue
		}

		userTask, err := database.GetUserTaskByUserIDTaskID(user.UserID, taskID)
		if err != nil {
			log.Printf("Failed to get user task: %v", err)
			continue
		}

		if userTask.Completed {
			totalUSDC += usdc
			continue
		}
		delete(senderMap, sender)
	}
	return senderMap, totalUSDC
}

func ParseAddress(address string) string {
	commonAddress := common.HexToAddress(address)
	return commonAddress.Hex()
}
