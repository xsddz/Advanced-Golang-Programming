package logger

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
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

	msg := fmt.Sprintf("[trace] [%v] [%v] [rows:%v] [%v] %v\n", traceID, float64(elapsed.Nanoseconds())/1e6, rows, err, sql)

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
func (l *Logger) Debug(ctx context.Context, message interface{}) (n int, err error) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprintf("[debug] [%v] %v\n", traceID, message)

	return l.Write([]byte(msg))
}

func (l *Logger) Critical(ctx context.Context, message interface{}) (n int, err error) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprintf("[critical] [%v] %v\n", traceID, message)

	return l.Write([]byte(msg))
}

func (l *Logger) Audit(ctx context.Context, message interface{}) (n int, err error) {
	traceID := ctx.Value("trace_id")

	msg := fmt.Sprintf("[audit] [%v] %v\n", traceID, message)

	return l.Write([]byte(msg))
}
