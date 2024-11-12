package database

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateCampaign(t *testing.T) {
	type args struct {
		name        string
		poolAddress string
		startTime   int64
		endTime     int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success - Create campaign",
			args: args{
				name:        "Test Campaign",
				poolAddress: "0x1234567890abcdef",
				startTime:   1633072800,
				endTime:     1633159200,
			},
			wantErr: false,
		},
		{
			name: "Fail - Empty pool address",
			args: args{
				name:        "Campaign with empty address",
				poolAddress: "",
				startTime:   1633072800,
				endTime:     1633159200,
			},
			wantErr: true,
		},
		{
			name: "Fail - Invalid start time",
			args: args{
				name:        "Campaign with invalid start time",
				poolAddress: "0x1234567890abcdef",
				startTime:   -1,
				endTime:     1633159200,
			},
			wantErr: true,
		},
		{
			name: "Fail - Invalid end time",
			args: args{
				name:        "Campaign with invalid end time",
				poolAddress: "0x1234567890abcdef",
				startTime:   1633072800,
				endTime:     0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := CreateCampaign(tt.args.name, tt.args.poolAddress, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCampaign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && gotID == 0 {
				t.Errorf("CreateCampaign() gotID = %v, expected valid ID", gotID)
			}

			if !tt.wantErr {
				var name, poolAddress string
				var startTime, endTime int64
				err := testDB.QueryRow("SELECT name, pool_address, start_time, end_time FROM campaigns WHERE pool_address = $1", tt.args.poolAddress).
					Scan(&name, &poolAddress, &startTime, &endTime)
				if err != nil {
					t.Errorf("Failed to retrieve campaign from database: %v", err)
				}
				if name != tt.args.name || poolAddress != tt.args.poolAddress || startTime != tt.args.startTime || endTime != tt.args.endTime {
					t.Errorf("Campaign data mismatch. Got name=%v, poolAddress=%v, startTime=%v, endTime=%v",
						name, poolAddress, startTime, endTime)
				}
			}
		})
	}
}

func TestGetCampaignByID(t *testing.T) {
	campaignID, _ := CreateCampaign("Test Campaign", "0x1234567890abcdef", 1633072800, 1633159200)
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Success - Get campaign by ID",
			args: args{
				id: campaignID,
			},
			wantErr: false,
		},
		{
			name: "Fail - Campaign not found",
			args: args{
				id: 9999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCampaignByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCampaignByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCampaignsByAddress(t *testing.T) {
	campaignID, _ := CreateCampaign("Test Campaign", "TestGetCampaignsByAddress", 1633072800, 1633159200)
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    []Campaign
		wantErr bool
	}{
		{
			name: "Success - Get campaigns by address",
			args: args{
				address: "TestGetCampaignsByAddress",
			},
			want: []Campaign{
				{
					CampaignID:  campaignID,
					Name:        "Test Campaign",
					PoolAddress: "TestGetCampaignsByAddress",
					StartTime:   1633072800,
					EndTime:     1633159200,
				},
			},
			wantErr: false,
		},
		{
			name: "Success - No campaigns found",
			args: args{
				address: "No campaigns found",
			},
			want:    []Campaign{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCampaignsByAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCampaignsByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.want) == 0 && len(got) == 0 {
				return
			}
			for i := range tt.want {
				if !reflect.DeepEqual(got[i].PoolAddress, tt.want[i].PoolAddress) {
					t.Errorf("GetCampaignsByAddress() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

func TestGetActiveCampaignAddresses(t *testing.T) {
	campaignID, _ := CreateCampaign("Test Campaign", "0x1234567890abcdef", time.Now().Unix()-100, time.Now().Unix()+100)
	campaign, _ := GetCampaignByID(campaignID)
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "Success - Get active campaign addresses",
			want:    []string{campaign.PoolAddress},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetActiveCampaignAddresses()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveCampaignAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveCampaignAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}
