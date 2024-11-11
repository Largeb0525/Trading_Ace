package server

import "time"

type CreateCampaignReq struct {
	Name                string  `json:"name" binding:"required"`
	PoolAddress         string  `json:"poolAddress" binding:"required"`
	StartAt             int64   `json:"startAt" binding:"required"`
	OnboardingReward    float64 `json:"onboardingReward" binding:"required"`
	OnboardingThreshold float64 `json:"onboardingThreshold" binding:"required"`
	PointPool           float64 `json:"pointPool" binding:"required"`
	Schedule            string  `json:"schedule" binding:"required"`
	Round               int64   `json:"round" binding:"required"`
}

type GetUserTaskStatusResp struct {
	Campaigns []CampaignResp `json:"campaigns"`
}

type CampaignResp struct {
	CampaignID  int              `json:"campaignId"`
	Name        string           `json:"name"`
	PoolAddress string           `json:"poolAddress"`
	StartTime   int64            `json:"startTime"`
	EndTime     int64            `json:"endTime"`
	Tasks       []TaskStatusResp `json:"tasks"`
}

type TaskStatusResp struct {
	TaskID      int     `json:"taskId"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Completed   bool    `json:"completed"`
	Amount      float64 `json:"amount"`
	Points      float64 `json:"points"`
	StartTime   int64   `json:"startTime"`
	EndTime     int64   `json:"endTime"`
}

type GetUserPointsHistoryResp struct {
	PointsHistory []PointsHistoryResp `json:"pointsHistory"`
	Total         float64             `json:"total"`
}

type PointsHistoryResp struct {
	CampaignID   int       `json:"campaignId"`
	CampaignName string    `json:"campaignName"`
	PoolAddress  string    `json:"poolAddress"`
	TaskID       int       `json:"taskId"`
	TaskType     string    `json:"taskType"`
	Description  string    `json:"description"`
	Points       float64   `json:"points"`
	Timestamp    time.Time `json:"timestamp"`
}
