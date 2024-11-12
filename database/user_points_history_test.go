package database

import (
	"reflect"
	"testing"
)

func TestCreateUserPointsHistory(t *testing.T) {
	userID, _ := CreateUser("TestCreateUserPointsHistory")
	campaignID, _ := CreateCampaign("TestCreateUserPointsHistory", "TestCreateUserPointsHistory", 0, 1)
	taskID, _ := CreateOnboardingTask(campaignID, "TestCreateUserPointsHistory", 0, 1, 0, 1)
	type args struct {
		userID     int
		taskID     int
		campaignID int
		points     float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Success - Create user points history",
			args: args{
				userID:     userID,
				taskID:     taskID,
				campaignID: campaignID,
				points:     0,
			},
			wantErr: false,
		},
		{
			name: "Fail - Create user points history",
			args: args{
				userID:     -1,
				taskID:     taskID,
				campaignID: campaignID,
				points:     0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateUserPointsHistory(tt.args.userID, tt.args.taskID, tt.args.campaignID, tt.args.points); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserPointsHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserPointsHistoryByUserID(t *testing.T) {
	userID, _ := CreateUser("TestGetUserPointsHistoryByUserID")
	CampaignID, _ := CreateCampaign("TestGetUserPointsHistoryByUserID", "TestGetUserPointsHistoryByUserID", 0, 1)
	TaskID, _ := CreateOnboardingTask(CampaignID, "TestGetUserPointsHistoryByUserID", 0, 1, 0, 1)
	// nolint
	CreateUserPointsHistory(userID, TaskID, CampaignID, 0)
	type args struct {
		userID int
	}
	tests := []struct {
		name    string
		args    args
		want    []UserPointsHistory
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Success - Get user points history by user ID",
			args: args{
				userID: userID,
			},
			want: []UserPointsHistory{
				{
					UserID:     userID,
					TaskID:     TaskID,
					CampaignID: CampaignID,
					Points:     0,
				},
			},
			wantErr: false,
		},
		{
			name: "Fail - Get user points history by user ID",
			args: args{
				userID: -1,
			},
			want:    []UserPointsHistory{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserPointsHistoryByUserID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserPointsHistoryByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i].Points, tt.want[i].Points) {
					t.Errorf("GetUserPointsHistoryByUserID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
