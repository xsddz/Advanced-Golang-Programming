package trace_test

import (
	"context"
	"os"
	"testing"
	"time"
	"yawebapp/library/infra/server"
	"yawebapp/library/infra/trace"
)

// go test -run "TestLogger"
func TestLogger(t *testing.T) {
	writelog := func(conf trace.LogConf) error {
		l, err := trace.NewLogger(conf)
		if err != nil {
			return err
		}

		l.AddWriter(os.Stdout)

		ctx := server.NewWebContextViaContext(context.TODO())

		fc := func() (string, int64) { return "select 1", 1 }

		l.Fatal(ctx, "fatallllll")
		l.Fatal(ctx, "fatallllll", 12345, "qwert", []int{1, 2, 3}, map[string]string{"q": "we", "a": "sd"})
		l.Error(ctx, "errorrrrrr")
		l.Error(ctx, "errorrrrrr", 12345, "qwert", []int{1, 2, 3}, map[string]string{"q": "we", "a": "sd"})
		l.Warn(ctx, "warnnnnnnnn")
		l.Warn(ctx, "warnnnnnnnn", 12345, "qwert", []int{1, 2, 3}, map[string]string{"q": "we", "a": "sd"})
		l.Info(ctx, "infoooooooo")
		l.Info(ctx, "infoooooooo", 12345, "qwert", []int{1, 2, 3}, map[string]string{"q": "we", "a": "sd"})
		l.Debug(ctx, "debugggggg")
		l.Debug(ctx, "debugggggg", 12345, "qwert", []int{1, 2, 3}, map[string]string{"q": "we", "a": "sd"})
		l.Trace(ctx, time.Now(), fc, err)

		return nil
	}

	cases := []struct {
		name string
		conf trace.LogConf
		want bool
	}{
		{"case 1", trace.LogConf{}, false},
		{"case 2", trace.LogConf{Level: trace.LOG_ALL}, false},
		{"case 3", trace.LogConf{Level: trace.LOG_ALL, HasColor: true}, false},
		{"case 4", trace.LogConf{Level: trace.LOG_ALL, HasColor: true, LogPath: "."}, false},
		{"case 5", trace.LogConf{Level: trace.LOG_ALL, HasColor: true, LogPath: ".", AppName: "trace"}, true},
		{"case 6", trace.LogConf{Level: trace.LOG_ALL, HasColor: false, LogPath: ".", AppName: "trace"}, true},
		{"case 7", trace.LogConf{Level: trace.LOG_ALL, HasColor: false, HasFileNum: true, LogPath: ".", AppName: "trace"}, true},
	}

	for _, c := range cases {
		err := writelog(c.conf)
		if got := (err == nil); got != c.want {
			t.Errorf("test %v failed. error => %v,", c.name, err)
		}
	}
}
