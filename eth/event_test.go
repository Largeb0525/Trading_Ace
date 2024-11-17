package eth

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/Largeb0525/Trading_Ace/database"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestParseAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				address: "0x0000000000000000000000000000000000000000",
			},
			want: "0x0000000000000000000000000000000000000000",
		},
		{
			name: "test2",
			args: args{
				address: "0x000000000000000000000000B4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
			},
			want: "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseAddress(tt.args.address); got != tt.want {
				t.Errorf("ParseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBlockTime(t *testing.T) {
	mockBlock := types.NewBlockWithHeader(&types.Header{
		Time: 1633072800,
	})

	patches := gomonkey.ApplyMethodFunc(client, "BlockByNumber", func(ctx context.Context, number *big.Int) (*types.Block, error) {
		if number == nil {
			return mockBlock, nil
		}
		switch number.Uint64() {
		case 100:
			return mockBlock, nil
		case 200:
			return nil, errors.New("block fetch error")
		default:
			return nil, errors.New("block not found")
		}
	})
	defer patches.Reset()

	tests := []struct {
		name        string
		blockNumber uint64
		wantTime    int64
		wantErr     bool
	}{
		{
			name:        "Success - valid block",
			blockNumber: 100,
			wantTime:    1633072800,
			wantErr:     false,
		},
		{
			name:        "Error - block fetch failed",
			blockNumber: 200,
			wantTime:    0,
			wantErr:     true,
		},
		{
			name:        "Error - block not found",
			blockNumber: 300,
			wantTime:    0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := getBlockTime(tt.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBlockTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantTime, gotTime)
		})
	}
}

func Test_timeToBlockNumber(t *testing.T) {
	mockBlock1 := types.NewBlockWithHeader(&types.Header{Time: 1633072800, Number: big.NewInt(1)})
	mockBlock2 := types.NewBlockWithHeader(&types.Header{Time: 1633076400, Number: big.NewInt(2)})
	mockBlock3 := types.NewBlockWithHeader(&types.Header{Time: 1633080000, Number: big.NewInt(3)})
	mockBlock4 := types.NewBlockWithHeader(&types.Header{Time: 1633083600, Number: big.NewInt(4)})

	patches := gomonkey.ApplyMethodFunc(client, "BlockByNumber", func(ctx context.Context, number *big.Int) (*types.Block, error) {
		if number == nil {
			return mockBlock4, nil
		}
		switch number.Uint64() {
		case 1:
			return mockBlock1, nil
		case 2:
			return mockBlock2, nil
		case 3:
			return mockBlock3, nil
		case 4:
			return mockBlock4, nil
		default:
			return nil, errors.New("block not found")
		}
	})
	defer patches.Reset()

	tests := []struct {
		name          string
		targetTime    int64
		afterBlock    uint64
		beforeBlock   uint64
		expectedBlock uint64
		expectError   bool
	}{
		{
			name:          "Exact match",
			targetTime:    1633076400,
			afterBlock:    0,
			beforeBlock:   0,
			expectedBlock: 2,
			expectError:   false,
		},
		{
			name:          "Before range",
			targetTime:    1633069200,
			afterBlock:    0,
			beforeBlock:   0,
			expectedBlock: 0,
			expectError:   false,
		},
		{
			name:          "After range",
			targetTime:    1633087200,
			afterBlock:    0,
			beforeBlock:   0,
			expectedBlock: 2,
			expectError:   false,
		},
		{
			name:          "Invalid block time",
			targetTime:    1633075000,
			afterBlock:    4,
			beforeBlock:   6,
			expectedBlock: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := timeToBlockNumber(tt.targetTime, tt.afterBlock, tt.beforeBlock)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBlock, block)
			}
		})
	}
}

func TestFetchSwapEvents(t *testing.T) {
	mockBlock1 := types.NewBlockWithHeader(&types.Header{Time: 1633072800, Number: big.NewInt(1)})
	mockBlock2 := types.NewBlockWithHeader(&types.Header{Time: 1633076400, Number: big.NewInt(2)})
	mockBlock3 := types.NewBlockWithHeader(&types.Header{Time: 1633080000, Number: big.NewInt(3)})
	mockBlock4 := types.NewBlockWithHeader(&types.Header{Time: 1633083600, Number: big.NewInt(4)})
	mockBlock6 := types.NewBlockWithHeader(&types.Header{Time: 1633086000, Number: big.NewInt(6)})

	patches := gomonkey.ApplyMethodFunc(client, "BlockByNumber", func(ctx context.Context, number *big.Int) (*types.Block, error) {
		if number == nil {
			return mockBlock6, nil
		}
		switch number.Uint64() {
		case 1:
			return mockBlock1, nil
		case 2:
			return mockBlock2, nil
		case 3:
			return mockBlock3, nil
		case 4:
			return mockBlock4, nil
		case 6:
			return mockBlock6, nil
		default:
			return nil, errors.New("block not found")
		}
	})
	defer patches.Reset()

	mockLogs := []types.Log{
		{
			Address: common.HexToAddress("0x123"),
			Topics:  []common.Hash{common.HexToHash(swapEventTopicHash)},
			Data:    []byte("mockData"),
		},
	}
	patches2 := gomonkey.ApplyMethodFunc(client, "FilterLogs", func(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
		if query.Addresses[0] == common.HexToAddress("0x123") {
			return mockLogs, nil
		}
		return nil, errors.New("no logs found")
	})
	defer patches2.Reset()

	tests := []struct {
		name        string
		poolAddress string
		startTime   int64
		endTime     int64
		expectLogs  []types.Log
		expectError bool
	}{
		{
			name:        "Success case",
			poolAddress: "0x123",
			startTime:   1633072800,
			endTime:     1633076400,
			expectLogs:  mockLogs,
			expectError: false,
		},
		{
			name:        "Block not found",
			poolAddress: "0x123",
			startTime:   1633083600,
			endTime:     1633086000,
			expectLogs:  nil,
			expectError: true,
		},
		{
			name:        "No logs found",
			poolAddress: "0x456",
			startTime:   1633072800,
			endTime:     1633076400,
			expectLogs:  nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := FetchSwapEvents(tt.poolAddress, tt.startTime, tt.endTime)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectLogs, logs)
			}
		})
	}
}

func TestParseSwapEvents(t *testing.T) {
	mockLogs := []types.Log{
		{
			Address: common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"),
			Topics: []common.Hash{
				common.HexToHash(swapEventTopicHash),
				common.HexToHash("0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"),
				common.HexToHash("0x19D1B048c8CDb4Cc280676627CE8c05756C5519e"),
			},
			Data:        []byte(`0x0000000000000000000000000000000000000000000000000000000017d784000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001c8d2577d58988e`),
			BlockNumber: 12345,
			TxHash:      common.HexToHash("0xfce9072349cf28089f7816c526fcf9d50ac277fc7ba9a43377228ea4a0604f70"),
		},
	}

	mockBlock := types.NewBlockWithHeader(&types.Header{Time: 1633083600, Number: big.NewInt(12345)})
	patches := gomonkey.ApplyMethodFunc(client, "BlockByNumber", func(ctx context.Context, number *big.Int) (*types.Block, error) {
		if number == nil {
			return mockBlock, nil
		}
		switch number.Uint64() {
		case 12345:
			return mockBlock, nil
		default:
			return nil, errors.New("block not found")
		}
	})
	defer patches.Reset()

	swapInfos, err := ParseSwapEvents(mockLogs)
	assert.NoError(t, err)
	assert.Len(t, swapInfos, 1)

	expectedSwapInfo := SwapInfo{
		Sender:      "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
		USDC:        3472328296227.68,
		Timestamp:   1633083600,
		PoolAddress: "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
		TxHash:      "0xfce9072349cf28089f7816c526fcf9d50ac277fc7ba9a43377228ea4a0604f70",
	}
	assert.Equal(t, expectedSwapInfo, swapInfos[0])
}

func Test_calculateTotalUSDC(t *testing.T) {

	patches := gomonkey.ApplyFunc(database.GetUserByAddress, func(address string) (*database.User, error) {
		if address == "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD" {
			return &database.User{UserID: 1}, nil
		} else if address == "0x19D1B048c8CDb4Cc280676627CE8c05756C5519e" {
			return &database.User{UserID: 3}, nil
		}
		return nil, nil
	})
	defer patches.Reset()

	patches2 := gomonkey.ApplyFunc(database.GetUserTaskByUserIDTaskID, func(userID, taskID int) (database.UserTask, error) {
		if userID == 1 {
			return database.UserTask{Completed: false}, nil
		} else if userID == 3 {
			return database.UserTask{Completed: true}, nil
		}
		return database.UserTask{}, nil
	})
	defer patches2.Reset()

	type args struct {
		senderMap map[string]float64
		taskID    int
		threshold float64
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]float64
		want1 float64
	}{
		{
			name: "test1",
			args: args{
				senderMap: map[string]float64{
					"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD": 100,
					"0x19D1B048c8CDb4Cc280676627CE8c05756C5519e": 200,
				},
				taskID:    1,
				threshold: 100,
			},
			want: map[string]float64{
				"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD": 100,
				"0x19D1B048c8CDb4Cc280676627CE8c05756C5519e": 200,
			},
			want1: 300,
		},
		{
			name: "test2",
			args: args{
				senderMap: map[string]float64{
					"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD": 100,
					"0x19D1B048c8CDb4Cc280676627CE8c05756C5519e": 200,
				},
				taskID:    2,
				threshold: 150,
			},
			want: map[string]float64{
				"0x19D1B048c8CDb4Cc280676627CE8c05756C5519e": 200,
			},
			want1: 200,
		},
		{
			name: "test3",
			args: args{
				senderMap: map[string]float64{
					"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD": 200,
					"0x19D1B048c8CDb4Cc280676627CE8c05756C5519e": 100,
				},
				taskID:    2,
				threshold: 150,
			},
			want: map[string]float64{
				"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD": 200,
				"0x19D1B048c8CDb4Cc280676627CE8c05756C5519e": 100,
			},
			want1: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := calculateTotalUSDC(tt.args.senderMap, tt.args.taskID, tt.args.threshold)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateTotalUSDC() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("calculateTotalUSDC() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
