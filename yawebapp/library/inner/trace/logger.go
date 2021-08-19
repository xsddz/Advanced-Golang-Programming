package trace

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
	"yawebapp/library/inner/helper"

	"gorm.io/gorm/logger"
)

var (
	logSourceFile string
)

func init() {
	_, logSourceFile, _, _ = runtime.Caller(0)
}

type Logger struct {
	WithFileNum bool // 输出日志中含调用logger的代码位置
}

func NewLogger() *Logger {
	return &Logger{
		WithFileNum: false,
	}
}

////////////////////////////////////////////////////////////////////////////////
// 实现 gorm logger.Interface 接口
// type Interface interface {
//     LogMode(LogLevel) Interface
//     Info(context.Context, string, ...interface{})
//     Warn(context.Context, string, ...interface{})
//     Error(context.Context, string, ...interface{})
//     Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
// }

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, message string, vals ...interface{}) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprint(vals...)
	msg = fmt.Sprintf("[info] [%v] %v, %v\n", traceID, message, msg)

	l.Write([]byte(msg))
}

func (l *Logger) Warn(ctx context.Context, message string, vals ...interface{}) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprint(vals...)
	msg = fmt.Sprintf("[warn] [%v] %v, %v\n", traceID, message, msg)

	l.Write([]byte(msg))
}

func (l *Logger) Error(ctx context.Context, message string, vals ...interface{}) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprint(vals...)
	msg = fmt.Sprintf("[error] [%v] %v, %v\n", traceID, message, msg)

	l.Write([]byte(msg))
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	traceID := ctx.Value("trace_id")
	elapsed := time.Since(begin)
	sql, rows := fc()

	colorFormat := "[trace] [%v] " + helper.YellowBold + "[%vms] " + helper.BlueBold + "[rows:%v] " + helper.MagentaBold + "%v" + helper.Reset + ": %v\n"
	msg := fmt.Sprintf(colorFormat, traceID, float64(elapsed.Nanoseconds())/1e6, rows, err, sql)

	l.Write([]byte(msg))
}

////////////////////////////////////////////////////////////////////////////////
// Write 实现 io.Writer 接口
// type Writer interface {
// 	Write(p []byte) (n int, err error)
// }

func (l *Logger) Write(p []byte) (n int, err error) {
	filenum := ""
	if l.WithFileNum {
		filenum = fileWithLineNum()
	}
	return fmt.Print(filenum, string(p))
}

func fileWithLineNum() string {
	for i := 0; i < 7; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if file == logSourceFile {
			continue
		}
		if strings.Index(file, "/pkg/") > 0 {
			continue
		}
		return file + ":" + strconv.FormatInt(int64(line), 10) + " "
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// 补充实现其他方法

func (l *Logger) Debug(ctx context.Context, message string, vals ...interface{}) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprint(vals...)
	msg = fmt.Sprintf("[debug] [%v] %v, %v\n", traceID, message, msg)

	l.Write([]byte(msg))
}

func (l *Logger) Critical(ctx context.Context, message string, vals ...interface{}) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprint(vals...)
	msg = fmt.Sprintf("[critical] [%v] %v, %v\n", traceID, message, msg)

	l.Write([]byte(msg))
}

func (l *Logger) Audit(ctx context.Context, message string, vals ...interface{}) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprint(vals...)
	msg = fmt.Sprintf("[audit] [%v] %v, %v\n", traceID, message, msg)

	l.Write([]byte(msg))
}
