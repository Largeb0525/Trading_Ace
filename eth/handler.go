package eth

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Largeb0525/Trading_Ace/database"
	"github.com/ethereum/go-ethereum/core/types"
)

func handleLogs(vLog types.Log) {
	swapInfos, err := ParseSwapEvents([]types.Log{vLog})
	if err != nil {
		log.Printf("Failed to parse Swap event: %v", err)
		return
	}

	for _, swapInfo := range swapInfos {
		insertSwapEventAndUpdateTask(swapInfo)
	}
}

func insertSwapEventAndUpdateTask(swapInfo SwapInfo) {
	userID, err := database.GetOrCreateUserID(swapInfo.Sender)
	if err != nil {
		log.Printf("Failed to get or create user ID: %v", err)
		return
	}

	err = database.InsertSwapEvent(userID, swapInfo.PoolAddress, swapInfo.USDC, swapInfo.Timestamp, swapInfo.TxHash)
	if err != nil {
		log.Printf("Failed to insert swap event: %v", err)
		return
	}

	fmt.Printf("Stored swap event: User: %s, Pool: %s, USDC: %f, Timestamp: %d\n", swapInfo.Sender, swapInfo.PoolAddress, swapInfo.USDC, swapInfo.Timestamp)
	go updateTask(userID, swapInfo.PoolAddress, swapInfo.USDC, swapInfo.Timestamp)
}

func updateTask(userID int, poolAddress string, usdc float64, t int64) {
	campaigns, err := database.GetCampaignsByAddress(poolAddress)
	if err != nil {
		log.Printf("Failed to get campaign: %v", err)
		return
	}

	for _, campaign := range campaigns {
		if t < campaign.StartTime || t > campaign.EndTime {
			return
		}
		processOnboardingTasks(campaign.CampaignID, userID, usdc)
		processSharePoolTask(campaign.CampaignID, userID, usdc, t)
	}
}

func insertUserTaskAndGetOnboardingUserTask(userID int, campaignID int, taskID int) (database.UserTask, error) {
	tasks, err := database.GetTasksByCampaignID(campaignID)
	if err != nil {
		return database.UserTask{}, err
	}

	onboardingUserTask := database.UserTask{}

	for _, task := range tasks {
		userTaskID, err := database.CreateUserTask(userID, task.TaskID, false, 0, 0)
		if err != nil {
			log.Printf("Failed to create user task: %v", err)
			continue
		}

		if task.Type == "onboarding" && taskID == task.TaskID {
			onboardingUserTask = database.UserTask{
				UserTaskID: userTaskID,
				UserID:     userID,
				TaskID:     taskID,
				Completed:  false,
				Amount:     0,
				Points:     0,
			}
		}
	}
	return onboardingUserTask, nil
}

func processOnboardingTasks(campaignID, userID int, usdc float64) {
	task, err := database.GetOnboardingTaskByCampaignID(campaignID)
	if err != nil {
		log.Printf("Failed to get onboarding tasks: %v", err)
		return
	}

	userTask, err := database.GetUserTaskByUserIDTaskID(userID, task.TaskID)
	if err == sql.ErrNoRows {
		userTask, err = insertUserTaskAndGetOnboardingUserTask(userID, campaignID, task.TaskID)
		if err != nil {
			log.Printf("Failed to insert user task: %v", err)
			return
		}
	} else if err != nil {
		log.Printf("Failed to get user task: %v", err)
		return
	}

	totalAmount := userTask.Amount + usdc
	if !userTask.Completed && totalAmount >= task.OnboardingThreshold {
		err = database.UpdateUserTask(userTask.UserTaskID, true, totalAmount, task.OnboardingReward)
		if err != nil {
			log.Printf("Failed to update user task: %v", err)
			return
		}
		err = database.CreateUserPointsHistory(userID, task.TaskID, campaignID, task.OnboardingReward)
		if err != nil {
			log.Printf("Failed to create user points history: %v", err)
			return
		}
	} else if !userTask.Completed {
		err = database.UpdateUserTask(userTask.UserTaskID, false, totalAmount, 0)
		if err != nil {
			log.Printf("Failed to update user task: %v", err)
			return
		}
	}
}

func processSharePoolTask(campaignID, userID int, USDC float64, t int64) {
	tasks, err := database.GetActiveTasksByCampaignID(campaignID, t)
	if err != nil {
		log.Printf("Failed to get active tasks: %v", err)
		return
	}
	for _, task := range tasks {
		if task.Type == "share_pool" {
			err = database.IncreaseUserTaskAmount(task.TaskID, userID, USDC)
			if err != nil {
				log.Printf("Failed to update task: %v", err)
				continue
			}
		}
	}
}
