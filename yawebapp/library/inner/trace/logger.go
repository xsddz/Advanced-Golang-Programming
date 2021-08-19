package trace

import (
	"context"
	"fmt"
	"time"
	"yawebapp/library/inner/helper"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
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

	colorFormat := "[trace] [%v] " + helper.YellowBold + "[%vms] " + helper.BlueBold + "[rows:%v] " + helper.RedBold + "%v " + helper.MagentaBold + "%v" + helper.Reset + ": %v\n"
	msg := fmt.Sprintf(colorFormat, traceID, float64(elapsed.Nanoseconds())/1e6, rows, utils.FileWithLineNum(), err, sql)

	l.Write([]byte(msg))
}

////////////////////////////////////////////////////////////////////////////////
// Write 实现 io.Writer 接口
// type Writer interface {
// 	Write(p []byte) (n int, err error)
// }

func (l *Logger) Write(p []byte) (n int, err error) {
	return fmt.Print(string(p))
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
