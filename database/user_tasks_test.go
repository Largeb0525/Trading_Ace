package database

import (
	"reflect"
	"testing"
)

func TestCreateUserTask(t *testing.T) {
	userID, _ := CreateUser("TestCreateUserTask")
	campaignID, _ := CreateCampaign("TestCreateUserTask", "TestCreateUserTask", 0, 1)
	taskID, _ := CreateOnboardingTask(campaignID, "TestCreateUserTask", 0, 1, 0, 1)
	type args struct {
		userID    int
		taskID    int
		completed bool
		amount    float64
		points    float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success - Create user task",
			args: args{
				userID:    userID,
				taskID:    taskID,
				completed: false,
				amount:    0,
				points:    0,
			},
			wantErr: false,
		},
		{
			name: "Fail - Create user task",
			args: args{
				userID:    -1,
				taskID:    taskID,
				completed: false,
				amount:    0,
				points:    0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateUserTask(tt.args.userID, tt.args.taskID, tt.args.completed, tt.args.amount, tt.args.points)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetUserTasksByUserID(t *testing.T) {
	userID, _ := CreateUser("TestGetUserTasksByUserID")
	campaignID, _ := CreateCampaign("TestGetUserTasksByUserID", "TestGetUserTasksByUserID", 0, 1)
	taskID, _ := CreateOnboardingTask(campaignID, "TestGetUserTasksByUserID", 0, 1, 0, 1)
	userTaskID, _ := CreateUserTask(userID, taskID, false, 0, 0)
	type args struct {
		userID int
	}
	tests := []struct {
		name    string
		args    args
		want    []UserTask
		wantErr bool
	}{
		{
			name: "Success - Get user tasks by user ID",
			args: args{
				userID: userID,
			},
			want: []UserTask{
				{
					UserTaskID: userTaskID,
					UserID:     userID,
					TaskID:     taskID,
					Completed:  false,
					Amount:     0,
					Points:     0,
				},
			},
			wantErr: false,
		},
		{
			name: "Fail - Get user tasks by user ID",
			args: args{
				userID: -1,
			},
			want:    []UserTask{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserTasksByUserID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserTasksByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i].UserTaskID, tt.want[i].UserTaskID) {
					t.Errorf("GetUserTasksByUserID() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}
