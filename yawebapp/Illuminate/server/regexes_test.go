package server_test

import (
	"errors"
	"testing"
	"yawebapp/Illuminate/server"
)

// go test -run "TestParseError"
func TestParseError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name        string
		args        args
		wantCode    int
		wantMessage string
	}{
		{
			name:        "case 1",
			args:        args{e: errors.New("[-1] ok 124214")},
			wantCode:    -1,
			wantMessage: "ok 124214",
		},
		{
			name:        "case 2",
			args:        args{e: errors.New("[0] ok 123")},
			wantCode:    0,
			wantMessage: "ok 123",
		},
		{
			name:        "case 3",
			args:        args{e: errors.New("[1]参数错误 123 ")},
			wantCode:    1,
			wantMessage: "参数错误 123 ",
		},
		{
			name:        "case 4",
			args:        args{e: errors.New("[ ] 啊啊啊")},
			wantCode:    -1,
			wantMessage: "[ ] 啊啊啊",
		},
		{
			name:        "case 5",
			args:        args{e: errors.New("[] 错误1224")},
			wantCode:    -1,
			wantMessage: "[] 错误1224",
		},
		{
			name:        "case 6",
			args:        args{e: errors.New("[200012] 参数错误1224")},
			wantCode:    200012,
			wantMessage: "参数错误1224",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotMessage := server.ParseError(tt.args.e)
			if gotCode != tt.wantCode {
				t.Errorf("ParseError() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
			if gotMessage != tt.wantMessage {
				t.Errorf("ParseError() gotMessage = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}
