package database

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	campaignID, _ := CreateCampaign("Test Campaign", "0x1234567890abcdef", 1633072800, 1633159200)
	type args struct {
		campaignID          int
		taskType            string
		description         string
		onboardingReward    float64
		onboardingThreshold float64
		pointsPool          float64
		startTime           int64
		endTime             int64
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Success - Create task",
			args: args{
				campaignID:          campaignID,
				taskType:            "create test",
				description:         "test",
				onboardingReward:    0,
				onboardingThreshold: 0,
				pointsPool:          0,
				startTime:           123,
				endTime:             234,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Fail - Empty task type",
			args: args{
				campaignID:          campaignID,
				taskType:            "",
				description:         "test",
				onboardingReward:    0,
				onboardingThreshold: 0,
				pointsPool:          0,
				startTime:           0,
				endTime:             0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Fail - Invalid campaign ID",
			args: args{
				campaignID:          -1,
				taskType:            "test",
				description:         "test",
				onboardingReward:    0,
				onboardingThreshold: 0,
				pointsPool:          0,
				startTime:           0,
				endTime:             0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Fail - Invalid start time",
			args: args{
				campaignID:          campaignID,
				taskType:            "test",
				description:         "test",
				onboardingReward:    0,
				onboardingThreshold: 0,
				pointsPool:          0,
				startTime:           -1,
				endTime:             0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Fail - Invalid end time",
			args: args{
				campaignID:          campaignID,
				taskType:            "test",
				description:         "test",
				onboardingReward:    0,
				onboardingThreshold: 0,
				pointsPool:          0,
				startTime:           0,
				endTime:             0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Fail - Invalid onboarding reward",
			args: args{
				campaignID:          campaignID,
				taskType:            "test",
				description:         "test",
				onboardingReward:    -1,
				onboardingThreshold: 0,
				pointsPool:          0,
				startTime:           0,
				endTime:             0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Fail - Invalid onboarding threshold",
			args: args{
				campaignID:          campaignID,
				taskType:            "test",
				description:         "test",
				onboardingReward:    0,
				onboardingThreshold: -1,
				pointsPool:          0,
				startTime:           0,
				endTime:             0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Fail - Invalid points pool",
			args: args{
				campaignID:          campaignID,
				taskType:            "test",
				description:         "test",
				onboardingReward:    0,
				onboardingThreshold: 0,
				pointsPool:          -1,
				startTime:           0,
				endTime:             0,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTask(tt.args.campaignID, tt.args.taskType, tt.args.description, tt.args.onboardingReward, tt.args.onboardingThreshold, tt.args.pointsPool, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var taskType, description string
				var onboardingReward, onboardingThreshold, pointsPool float64
				var startTime, endTime int64
				query := `SELECT type, description, onboarding_reward, onboarding_threshold, points_pool, start_time, end_time FROM tasks WHERE task_id = $1`
				err := db.QueryRow(query, got).Scan(&taskType, &description, &onboardingReward, &onboardingThreshold, &pointsPool, &startTime, &endTime)
				if err != nil {
					t.Errorf("CreateTask() error = %v", err)
					return
				}
			}
		})
	}
}

func TestGetTasksByTaskIDs(t *testing.T) {
	campaignID, _ := CreateCampaign("test", "test", 0, 0)
	task1, _ := CreateTask(campaignID, "test", "test", 0, 0, 0, 0, 0)
	task2, _ := CreateTask(campaignID, "test", "test", 0, 0, 0, 0, 0)
	type args struct {
		taskIDs []int
	}
	tests := []struct {
		name    string
		args    args
		want    []Task
		wantErr bool
	}{
		{
			name: "Success - Get tasks by task IDs",
			args: args{
				taskIDs: []int{task1, task2},
			},
			wantErr: false,
		},
		{
			name: "Fail - Empty task IDs",
			args: args{
				taskIDs: []int{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTasksByTaskIDs(tt.args.taskIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTasksByTaskIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 && len(tt.args.taskIDs) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTasksByTaskIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateOnboardingTask(t *testing.T) {
	campaignID, _ := CreateCampaign("test", "test", 0, 1)
	type args struct {
		campaignID          int
		description         string
		onboardingReward    float64
		onboardingThreshold float64
		startTime           int64
		endTime             int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success - Create onboarding task",
			args: args{
				campaignID:          campaignID,
				description:         "test",
				onboardingReward:    0,
				onboardingThreshold: 0,
				startTime:           123,
				endTime:             456,
			},
			want:    "onboarding",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateOnboardingTask(tt.args.campaignID, tt.args.description, tt.args.onboardingReward, tt.args.onboardingThreshold, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOnboardingTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tasks, _ := GetTasksByTaskIDs([]int{got})
			if tasks[0].Type != tt.want {
				t.Errorf("CreateOnboardingTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateSharePoolTask(t *testing.T) {
	campaignID, _ := CreateCampaign("test", "test", 0, 1)
	type args struct {
		campaignID  int
		description string
		pointsPool  float64
		startTime   int64
		endTime     int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success - Create share pool task",
			args: args{
				campaignID:  campaignID,
				description: "test",
				pointsPool:  0,
				startTime:   123,
				endTime:     456,
			},
			want:    "share_pool",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSharePoolTask(tt.args.campaignID, tt.args.description, tt.args.pointsPool, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSharePoolTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tasks, _ := GetTasksByTaskIDs([]int{got})
			if tasks[0].Type != tt.want {
				t.Errorf("CreateSharePoolTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTasksByCampaignID(t *testing.T) {
	campaignID, _ := CreateCampaign("test", "TestGetTasksByCampaignID", 0, 1)
	campaign, _ := GetCampaignByID(campaignID)
	taskID, _ := CreateOnboardingTask(campaign.CampaignID, "test", 0, 0, 123, 456)
	type args struct {
		campaignID int
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "Success - Get tasks by campaign ID",
			args: args{
				campaignID: campaign.CampaignID,
			},
			want:    []int{taskID},
			wantErr: false,
		},
		{
			name: "Success - Get no tasks by campaign ID",
			args: args{
				campaignID: -1,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTasksByCampaignID(tt.args.campaignID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTasksByCampaignID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := 0; i < len(got); i++ {
				if got[i].TaskID != tt.want[i] {
					t.Errorf("GetTasksByCampaignID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetActiveTasksByCampaignID(t *testing.T) {
	campaignID, _ := CreateCampaign("test", "TestGetActiveTasksByCampaignID", 456, 789)
	campaign, _ := GetCampaignByID(campaignID)
	taskID, _ := CreateOnboardingTask(campaign.CampaignID, "test", 0, 0, 456, 789)
	type args struct {
		campaignID int
		timestamp  int64
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "Success - Get active tasks by campaign ID",
			args: args{
				campaignID: campaign.CampaignID,
				timestamp:  500,
			},
			want:    []int{taskID},
			wantErr: false,
		},
		{
			name: "Fail - No active tasks",
			args: args{
				campaignID: campaign.CampaignID,
				timestamp:  1000,
			},
			want:    []int{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetActiveTasksByCampaignID(tt.args.campaignID, tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveTasksByCampaignID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := 0; i < len(got); i++ {
				if got[i].TaskID != tt.want[i] {
					t.Errorf("GetActiveTasksByCampaignID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetOnboardingTaskByCampaignID(t *testing.T) {
	campaignID, _ := CreateCampaign("test", "TestGetOnboardingTaskByCampaignID", 456, 789)
	campaign, _ := GetCampaignByID(campaignID)
	taskID, _ := CreateOnboardingTask(campaign.CampaignID, "test", 0, 0, 456, 789)
	type args struct {
		campaignID int
	}
	tests := []struct {
		name    string
		args    args
		want    *Task
		wantErr bool
	}{
		{
			name: "Success - Get onboarding task by campaign ID",
			args: args{
				campaignID: campaign.CampaignID,
			},
			want: &Task{
				TaskID:      taskID,
				CampaignID:  campaign.CampaignID,
				Type:        "onboarding",
				Description: "test",
				StartTime:   456,
				EndTime:     789,
			},
			wantErr: false,
		},
		{
			name: "Fail - Get onboarding task by campaign ID",
			args: args{
				campaignID: -1,
			},
			want:    &Task{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOnboardingTaskByCampaignID(tt.args.campaignID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOnboardingTaskByCampaignID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOnboardingTaskByCampaignID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetExpiredSharePoolTasks(t *testing.T) {
	now := time.Now().Unix()
	campaignID, _ := CreateCampaign("test", "TestGetExpiredSharePoolTasks", now-1000, now-500)
	campaign, _ := GetCampaignByID(campaignID)
	taskID, _ := CreateSharePoolTask(campaign.CampaignID, "test", 0, now-1000, now-500)
	type args struct {
		now           int64
		lastCheckTime int64
	}
	tests := []struct {
		name    string
		args    args
		want    []Task
		wantErr bool
	}{
		{
			name: "Success - Get expired share pool tasks",
			args: args{
				now:           now,
				lastCheckTime: now - 2000,
			},
			want: []Task{
				{
					TaskID:      taskID,
					CampaignID:  campaign.CampaignID,
					Type:        "share_pool",
					Description: "test",
					StartTime:   now - 1000,
					EndTime:     now - 500,
				},
			},
			wantErr: false,
		},
		{
			name: "Fail - No expired share pool tasks",
			args: args{
				now:           now,
				lastCheckTime: now,
			},
			want:    []Task{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetExpiredSharePoolTasks(tt.args.now, tt.args.lastCheckTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExpiredSharePoolTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExpiredSharePoolTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}
