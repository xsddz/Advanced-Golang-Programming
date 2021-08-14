package exponentialbackoff_test

import (
	"demo/exponentialbackoff"
	"fmt"
	"testing"
)

// TestConnectRetry -
// run:
//     go test -v -run="TestConnectRetry"
func TestConnectRetry(t *testing.T) {
	type args struct {
		network string
		address string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case 1",
			args: args{
				network: "tcp",
				address: "192.168.0.1:8078",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := exponentialbackoff.ConnectRetry(tt.args.network, tt.args.address)
			fmt.Println("conn:", conn, err)
		})
	}
}
