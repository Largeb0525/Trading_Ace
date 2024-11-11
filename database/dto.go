package database

type User struct {
	UserID    int
	Address   string
	CreatedAt int64
}
type Campaign struct {
	CampaignID  int
	Name        string
	PoolAddress string
	StartTime   int64
	EndTime     int64
	CreatedAt   int64
	UpdatedAt   int64
}

type Task struct {
	TaskID              int
	CampaignID          int
	Type                string
	Description         string
	OnboardingReward    float64
	OnboardingThreshold float64
	PointsPool          float64
	StartTime           int64
	EndTime             int64
	CreatedAt           int64
	UpdatedAt           int64
}

type UserTask struct {
	UserTaskID int
	UserID     int
	TaskID     int
	Completed  bool
	Amount     float64
	Points     float64
	CreatedAt  int64
	UpdatedAt  int64
}

type UserPointsHistory struct {
	HistoryID  int
	UserID     int
	TaskID     int
	CampaignID int
	Points     float64
	CreatedAt  int64
}

type UserSwap struct {
	SwapID          int
	UserID          int
	TransactionHash string
	PoolAddress     string
	AmountUSDC      float64
	AmountWETH      float64
	SwapTime        int64
	CreatedAt       int64
}
