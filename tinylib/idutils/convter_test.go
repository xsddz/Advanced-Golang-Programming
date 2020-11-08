package idutils_test

import (
	"Advanced-Golang-Programming/tinylib/idutils"
	"testing"
)

// TestA34Encode -
// run:
//     go test -v -run="TestA34Encode"
func TestA34Encode(t *testing.T) {
	type args struct {
		num  int64
		salt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				num:  3876269,
				salt: "",
			},
			want: "3wn6x",
		},
		{
			name: "case 2",
			args: args{
				num:  3876269,
				salt: "abc1",
			},
			want: "56r9p",
		},
		{
			name: "case 3",
			args: args{
				num:  3876269,
				salt: "abc2",
			},
			want: "c89dp",
		},
		{
			name: "case 4",
			args: args{
				num:  4000000,
				salt: "abc2",
			},
			want: "cvekc",
		},
		{
			name: "case 5",
			args: args{
				num:  4000001,
				salt: "abc2",
			},
			want: "cvekf",
		},
		{
			name: "case 6",
			args: args{
				num:  4000002,
				salt: "abc2",
			},
			want: "cvek6",
		},
		{
			name: "case 7",
			args: args{
				num:  9040764591583297537,
				salt: "abc2",
			},
			want: "febac2d78mbwb",
		},
		{
			name: "case 8",
			args: args{
				num:  9040764591583297538,
				salt: "abc2",
			},
			want: "febac2d78mbwo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idutils.A34Encode(tt.args.num, tt.args.salt); got != tt.want {
				t.Errorf("A34Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestA34Decode -
// run:
//     go test -v -run="TestA34Decode"
func TestA34Decode(t *testing.T) {
	type args struct {
		input string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "case 1",
			args: args{
				input: "3wn6x",
				salt:  "",
			},
			want: 3876269,
		},
		{
			name: "case 2",
			args: args{
				input: "56r9p",
				salt:  "abc1",
			},
			want: 3876269,
		},
		{
			name: "case 3",
			args: args{
				input: "c89dp",
				salt:  "abc2",
			},
			want: 3876269,
		},
		{
			name: "case 4",
			args: args{
				input: "cvekc",
				salt:  "abc2",
			},
			want: 4000000,
		},
		{
			name: "case 5",
			args: args{
				input: "cvekf",
				salt:  "abc2",
			},
			want: 4000001,
		},
		{
			name: "case 6",
			args: args{
				input: "cvek6",
				salt:  "abc2",
			},
			want: 4000002,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idutils.A34Decode(tt.args.input, tt.args.salt); got != tt.want {
				t.Errorf("A34Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestA34Decode -
// run:
//     go test -v -run="TestA34Decode"
func TestA58Encode(t *testing.T) {
	type args struct {
		num  int64
		salt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				num:  3876269,
				salt: "",
			},
			want: "LsHE",
		},
		{
			name: "case 2",
			args: args{
				num:  3876269,
				salt: "abc1",
			},
			want: "Nfgb",
		},
		{
			name: "case 3",
			args: args{
				num:  3876269,
				salt: "abc2",
			},
			want: "fimr",
		},
		{
			name: "case 4",
			args: args{
				num:  4000000,
				salt: "abc2",
			},
			want: "YN6j",
		},
		{
			name: "case 5",
			args: args{
				num:  4000001,
				salt: "abc2",
			},
			want: "YN6x",
		},
		{
			name: "case 6",
			args: args{
				num:  4000002,
				salt: "abc2",
			},
			want: "YN6q",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idutils.A58Encode(tt.args.num, tt.args.salt); got != tt.want {
				t.Errorf("A58Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestA58Decode -
// run:
//     go test -v -run="TestA58Decode"
func TestA58Decode(t *testing.T) {
	type args struct {
		input string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "case 1",
			args: args{
				input: "LsHE",
				salt:  "",
			},
			want: 3876269,
		},
		{
			name: "case 2",
			args: args{
				input: "Nfgb",
				salt:  "abc1",
			},
			want: 3876269,
		},
		{
			name: "case 3",
			args: args{
				input: "fimr",
				salt:  "abc2",
			},
			want: 3876269,
		},
		{
			name: "case 4",
			args: args{
				input: "YN6j",
				salt:  "abc2",
			},
			want: 4000000,
		},
		{
			name: "case 5",
			args: args{
				input: "YN6x",
				salt:  "abc2",
			},
			want: 4000001,
		},
		{
			name: "case 6",
			args: args{
				input: "YN6q",
				salt:  "abc2",
			},
			want: 4000002,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idutils.A58Decode(tt.args.input, tt.args.salt); got != tt.want {
				t.Errorf("A58Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
