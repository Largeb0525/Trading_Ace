package database

import (
	"testing"
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
				err := testDB.QueryRow("SELECT name, pool_address, start_time, end_time FROM campaigns WHERE campaign_id = $1", gotID).
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
