package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Largeb0525/Trading_Ace/database"
	"github.com/Largeb0525/Trading_Ace/eth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func CreateCampaignHandler(c *gin.Context) {
	var req CreateCampaignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	startTime := req.StartAt
	duration, err := time.ParseDuration(req.Schedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration format"})
		return
	}
	endTime := startTime + int64(duration.Seconds())*req.Round

	campaignID, err := database.CreateCampaign(req.Name, req.PoolAddress, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return
	}

	_, err = database.CreateOnboardingTask(campaignID, "", req.OnboardingReward, req.OnboardingThreshold, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create onboarding task"})
		return
	}

	for i := 0; i < int(req.Round); i++ {
		endTime = startTime + int64(duration.Seconds())
		describe := fmt.Sprintf("Round %d", i+1)
		_, err = database.CreateSharePoolTask(campaignID, describe, req.PointPool, startTime, endTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share pool task"})
			return
		}
		startTime = endTime
	}

	eth.AddAddresses([]string{req.PoolAddress})

	c.JSON(http.StatusOK, gin.H{"message": "Campaign and tasks created successfully"})
}

func GetUserTaskStatusHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("userID"))
	inputAddress := c.Query("userAddress")

	if userID == 0 && inputAddress != "" {
		userAddress := eth.ParseAddress(c.Query("userAddress"))
		user, err := database.GetUserByAddress(userAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user by address"})
			return
		}
		userID = user.UserID
	} else if userID == 0 && inputAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	userTasks, err := database.GetUserTasksByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tasks"})
		return
	}
	resp, err := buildUserTaskStatusResponse(userTasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetUserPointsHistoryHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("userID"))
	inputAddress := c.Query("userAddress")

	if userID == 0 && inputAddress != "" {
		userAddress := eth.ParseAddress(c.Query("userAddress"))
		user, err := database.GetUserByAddress(userAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user by address"})
			return
		}
		userID = user.UserID
	} else if userID == 0 && inputAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	UserPointsHistory, err := database.GetUserPointsHistoryByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user points history"})
		return
	}

	resp, err := buildUserPointsHistoryResponse(UserPointsHistory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func buildUserTaskStatusResponse(userTasks []database.UserTask) (GetUserTaskStatusResp, error) {
	userTasksMap := make(map[int]database.UserTask)
	taskIDs := make([]int, len(userTasks))
	for i, userTask := range userTasks {
		userTasksMap[userTask.TaskID] = userTask
		taskIDs[i] = userTask.TaskID
	}

	tasks, err := database.GetTasksByTaskIDs(taskIDs)
	if err != nil {
		return GetUserTaskStatusResp{}, fmt.Errorf("Failed to get tasks by IDs: %w", err)
	}

	campaignRespMap := make(map[int]CampaignResp)
	for _, task := range tasks {
		campaignResp, exists := campaignRespMap[task.CampaignID]
		if !exists {
			campaign, err := database.GetCampaignByID(task.CampaignID)
			if err != nil {
				return GetUserTaskStatusResp{}, fmt.Errorf("Failed to get campaign: %w", err)
			}
			campaignResp = CampaignResp{
				CampaignID:  campaign.CampaignID,
				Name:        campaign.Name,
				PoolAddress: campaign.PoolAddress,
				StartTime:   campaign.StartTime,
				EndTime:     campaign.EndTime,
				Tasks:       []TaskStatusResp{},
			}
		}

		userTask := userTasksMap[task.TaskID]
		taskResp := TaskStatusResp{
			TaskID:      task.TaskID,
			Type:        task.Type,
			Description: task.Description,
			Completed:   userTask.Completed,
			Amount:      userTask.Amount,
			Points:      userTask.Points,
			StartTime:   task.StartTime,
			EndTime:     task.EndTime,
		}

		campaignResp.Tasks = append(campaignResp.Tasks, taskResp)
		campaignRespMap[task.CampaignID] = campaignResp
	}

	var campaigns []CampaignResp
	for _, campaign := range campaignRespMap {
		campaigns = append(campaigns, campaign)
	}

	return GetUserTaskStatusResp{Campaigns: campaigns}, nil
}

func buildUserPointsHistoryResponse(UserPointsHistory []database.UserPointsHistory) (GetUserPointsHistoryResp, error) {
	total := 0.0
	campaignMap := make(map[int]database.Campaign)
	pointsHistory := []PointsHistoryResp{}
	for _, history := range UserPointsHistory {
		campaign, ok := campaignMap[history.CampaignID]
		if !ok {
			c, err := database.GetCampaignByID(history.CampaignID)
			if err != nil {
				return GetUserPointsHistoryResp{}, fmt.Errorf("Failed to get campaign: %w", err)
			}
			campaign = *c
			campaignMap[history.CampaignID] = campaign
		}

		tasks, err := database.GetTasksByTaskIDs([]int{history.TaskID})
		if err != nil {
			return GetUserPointsHistoryResp{}, fmt.Errorf("Failed to get tasks by IDs: %w", err)
		}

		task := tasks[0]

		pointsHistoryResp := PointsHistoryResp{
			CampaignID:   history.CampaignID,
			CampaignName: campaign.Name,
			PoolAddress:  campaign.PoolAddress,
			TaskID:       history.TaskID,
			TaskType:     task.Type,
			Description:  task.Description,
			Points:       history.Points,
			Timestamp:    time.Unix(history.CreatedAt, 0),
		}
		pointsHistory = append(pointsHistory, pointsHistoryResp)

		total += history.Points
	}

	return GetUserPointsHistoryResp{PointsHistory: pointsHistory, Total: total}, nil
}

func ProcessSharePoolTicker() {
	duration, err := time.ParseDuration(viper.GetString("server.ticker"))
	if err != nil {
		log.Printf("Failed to parse duration: %v", err)
		return
	}
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		checkAndProcessSharePoolTasks(duration)
	}
}

func checkAndProcessSharePoolTasks(duration time.Duration) {
	now := time.Now().Unix()
	lastChecked := now - int64(duration.Seconds())
	tasks, err := database.GetExpiredSharePoolTasks(now, lastChecked)
	if err != nil {
		log.Printf("Failed to retrieve expired share pool tasks: %v", err)
		return
	}
	OnboardingTaskIDMap := make(map[int]database.Task)
	campaignMap := make(map[int]database.Campaign)

	for _, task := range tasks {
		campaign, ok := campaignMap[task.CampaignID]
		if !ok {
			newCampaign, err := database.GetCampaignByID(task.CampaignID)
			if err != nil {
				log.Printf("Failed to retrieve campaign for task %d: %v", task.TaskID, err)
				continue
			}
			campaign = *newCampaign
			campaignMap[campaign.CampaignID] = campaign
		}
		_, ok = OnboardingTaskIDMap[campaign.CampaignID]
		if !ok {
			onboardingTask, err := database.GetOnboardingTaskByCampaignID(campaign.CampaignID)
			if err != nil {
				log.Printf("Failed to retrieve onboarding tasks for campaign %d: %v", campaign.CampaignID, err)
				continue
			}
			OnboardingTaskIDMap[onboardingTask.CampaignID] = *onboardingTask
		}

		swapEvents, err := eth.FetchSwapEvents(campaign.PoolAddress, task.StartTime, task.EndTime)
		if err != nil {
			log.Printf("Failed to fetch swap events for task %d: %v", task.TaskID, err)
			continue
		}
		swapInfos, err := eth.ParseSwapEvents(swapEvents)
		if err != nil {
			log.Printf("Failed to parse swap events for task %d: %v", task.TaskID, err)
			continue
		}
		eth.ProcessSwapInfos(task, swapInfos, OnboardingTaskIDMap[campaign.CampaignID])
	}
}
