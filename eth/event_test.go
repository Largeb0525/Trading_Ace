package eth

import (
	"testing"
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
