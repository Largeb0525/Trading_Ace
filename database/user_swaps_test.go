package database

import "testing"

func TestInsertSwapEvent(t *testing.T) {
	userID, _ := CreateUser("TestInsertSwapEvent")
	type args struct {
		userID      int
		poolAddress string
		usdc        float64
		swapTime    int64
		txHash      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Success - Insert swap event",
			args: args{
				userID:      userID,
				poolAddress: "0x123",
				usdc:        0,
				swapTime:    123,
				txHash:      "0x123",
			},
			wantErr: false,
		},
		{
			name: "Fail - Insert swap event",
			args: args{
				userID:      -1,
				poolAddress: "0x123",
				usdc:        0,
				swapTime:    123,
				txHash:      "0x123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertSwapEvent(tt.args.userID, tt.args.poolAddress, tt.args.usdc, tt.args.swapTime, tt.args.txHash); (err != nil) != tt.wantErr {
				t.Errorf("InsertSwapEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
