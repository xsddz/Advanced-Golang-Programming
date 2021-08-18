package server

import (
	"fmt"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

// Write 实现 io.Writer 接口
func (l *Logger) Write(p []byte) (n int, err error) {
	return fmt.Print(string(p))
}

func (l *Logger) Debug(ctx WebContext, message interface{}) (n int, err error) {
	traceID, _ := ctx.Get("trace_id")
	msg := fmt.Sprintf("[debug] [%v] %v\n", traceID, message)
	return l.Write([]byte(msg))
}

func (l *Logger) Info(ctx WebContext, message interface{}) (n int, err error) {
	traceID, _ := ctx.Get("trace_id")
	msg := fmt.Sprintf("[info] [%v] %v\n", traceID, message)
	return l.Write([]byte(msg))
}

func (l *Logger) Error(ctx WebContext, message interface{}) (n int, err error) {
	traceID, _ := ctx.Get("trace_id")
	msg := fmt.Sprintf("[error] [%v] %v\n", traceID, message)
	return l.Write([]byte(msg))
}

func (l *Logger) Critical(ctx WebContext, message interface{}) (n int, err error) {
	traceID, _ := ctx.Get("trace_id")
	msg := fmt.Sprintf("[critical] [%v] %v\n", traceID, message)
	return l.Write([]byte(msg))
}

func (l *Logger) Audit(ctx WebContext, message interface{}) (n int, err error) {
	traceID, _ := ctx.Get("trace_id")
	msg := fmt.Sprintf("[audit] [%v] %v\n", traceID, message)
	return l.Write([]byte(msg))
}
