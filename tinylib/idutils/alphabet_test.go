package idutils_test

import (
	"Advanced-Golang-Programming/tinylib/idutils"
	"testing"
)

// TestAlphabetShuffle -
// run:
//     go test -v -run="TestAlphabetShuffle"
func TestAlphabetShuffle(t *testing.T) {
	type args struct {
		alphabet string
		salt     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case alphabet 34 with no salt",
			args: args{
				alphabet: idutils.Alpha34,
				salt:     "",
			},
			want: "123456789abcdefghijkmnopqrstuvwxyz",
		},
		{
			name: "case alphabet 34 with salt: abc1",
			args: args{
				alphabet: idutils.Alpha34,
				salt:     "abc1",
			},
			want: "2e5y89zfshxtucqdik3jor41mwbang6p7v",
		},
		{
			name: "case alphabet 34 with salt: abc2",
			args: args{
				alphabet: idutils.Alpha34,
				salt:     "abc2",
			},
			want: "51cf6dgkznjst4a2imqxy9r3uweboh8p7v",
		},
		{
			name: "case alphabet 58 with no salt",
			args: args{
				alphabet: idutils.Alpha58,
				salt:     "",
			},
			want: "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz",
		},
		{
			name: "case alphabet 58 with salt: abc1",
			args: args{
				alphabet: idutils.Alpha58,
				salt:     "abc1",
			},
			want: "YA7eqVjtSZKcXbwrgmaN2kUChudGvyo5Hi3xWPLDTns649z1MpfRE8JBFQ",
		},
		{
			name: "case alphabet 58 with salt: abc2",
			args: args{
				alphabet: idutils.Alpha58,
				salt:     "abc2",
			},
			want: "5JR6dk8ngpTscrKMmEhfYV4weWzZaNjxquAGbotHXD2v7CU3P1iSy9LBFQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idutils.AlphabetShuffle(tt.args.alphabet, tt.args.salt); got != tt.want {
				t.Errorf("AlphabetShuffle() = %v, want %v", got, tt.want)
			}
		})
	}
}
