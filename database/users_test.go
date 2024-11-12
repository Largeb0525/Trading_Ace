package database

import (
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success - Create user",
			args: args{
				address: "0x123",
			},
			wantErr: false,
		},
		{
			name: "Fail - Create user",
			args: args{
				address: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateUser(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetUserByAddress(t *testing.T) {
	userID, _ := CreateUser("TestGetUserByAddress")
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "Success - Get user by address",
			args: args{
				address: "TestGetUserByAddress",
			},
			want: &User{
				UserID:  userID,
				Address: "TestGetUserByAddress",
			},
			wantErr: false,
		},
		{
			name: "Fail - Get user by address",
			args: args{
				address: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserByAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}
			if !reflect.DeepEqual(got.Address, tt.want.Address) {
				t.Errorf("GetUserByAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOrCreateUserID(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success - Get or create user ID",
			args: args{
				address: "TestGetOrCreateUserID",
			},
			wantErr: false,
		},
		{
			name: "Fail - Get or create user ID",
			args: args{
				address: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetOrCreateUserID(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrCreateUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
