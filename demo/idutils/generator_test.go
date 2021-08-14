package idutils_test

import (
	"demo/idutils"
	"fmt"
	"sync"
	"testing"
	"time"
)

// BenchmarkDefaultNextIDParallel -
// run:
//     go test -v -run="none" -bench="BenchmarkDefaultNextIDParallel" -benchtime="1s" > idddddd.log
//     将输出的id重定向到文件中，进行分析，是否有重复的id
//     sort idddddd.log | uniq -c | sort -k1,1 -nr | head
func BenchmarkDefaultNextIDParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := idutils.DefaultNextID()
			fmt.Println("id =>", id) //, idutils.A58Encode(id, "abc2"), idutils.A34Encode(id, "abc2"))
		}
	})
}

// TestDefaultNextID -
// run:
//     go test -v -run="TestDefaultNextID"
func TestDefaultNextID(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		{
			name: "case two goroutine",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(2)
			var got1, got2 int64
			go func() {
				defer wg.Done()
				got1 = idutils.DefaultNextID()
			}()
			go func() {
				defer wg.Done()
				got2 = idutils.DefaultNextID()
			}()
			wg.Wait()
			if got1 == got2 {
				t.Errorf("DefaultNextID() = %v, %v are equal", got1, got2)
			}
		})
	}
}

// TestIDGenerator -
// run:
//     go test -v -run="TestIDGenerator"
func TestIDGenerator(t *testing.T) {
	type args struct {
		currTime time.Time
		shardID  uint64
		incrID   uint64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "case 1",
			args: args{
				currTime: time.Date(2020, time.August, 20, 17, 06, 23, 0, time.UTC),
				shardID:  3876269,
				incrID:   1,
			},
			want: 56903100827123713,
		},
		{
			name: "case 2",
			args: args{
				currTime: time.Date(2020, time.August, 20, 17, 06, 23, 0, time.UTC),
				shardID:  3876269,
				incrID:   2,
			},
			want: 56903100827123714,
		},
		{
			name: "case 3",
			args: args{
				currTime: time.Date(2020, time.August, 20, 17, 06, 23, 0, time.UTC),
				shardID:  4000002,
				incrID:   2,
			},
			want: 56903100824322050,
		},
		{
			name: "case 4",
			args: args{
				currTime: time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
				shardID:  4000002,
				incrID:   2,
			},
			want: 9040764591583297538,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idutils.IDGenerator(tt.args.currTime, tt.args.shardID, tt.args.incrID); got != tt.want {
				t.Errorf("Generator() = %v, want %v", got, tt.want)
			}
		})
	}
}
