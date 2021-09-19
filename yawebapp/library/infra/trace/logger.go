package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"yawebapp/library/infra/helper"

	"gorm.io/gorm/logger"
)

var logSourceFile = "/trace/logger.go"

// ReportHandler 告警处理器标准
type ReportHandler func(context.Context, string, interface{}) error

type LogLevel int

const (
	// 如果设置优先级为WARN，那么OFF、FATAL、ERROR、WARN 4个级别的log能正常输出，而INFO、DEBUG、TRACE、ALL级别的log则会被忽略
	LOG_OFF   LogLevel = iota // 0
	LOG_FATAL                 // 1
	LOG_ERROR                 // 2
	LOG_WARN                  // 3
	LOG_INFO                  // 4
	LOG_DEBUG                 // 5
	LOG_TRACE                 // 6
	LOG_ALL                   // 7
)

// LogConf 日志配置
type LogConf struct {
	Level LogLevel

	HasColor   bool // 输出日志中带颜色（在终端下有效）
	HasFileNum bool // 输出日志中含调用logger的代码位置

	LogPath string
	AppName string
}

// Logger -
type Logger struct {
	conf       *LogConf
	wr         io.Writer
	reportFunc ReportHandler // fatal时调用的告警处理器
}

// NewLogger -
func NewLogger(conf LogConf) (*Logger, error) {
	l := &Logger{
		conf: &conf,
		wr:   os.Stdout,
	}

	if conf.LogPath == "" || conf.AppName == "" {
		j, _ := json.Marshal(conf)
		err := fmt.Errorf("invalid log conf: %s", j)
		l.Fatal(context.TODO(), err.Error())
		return nil, err
	}

	logFile := fmt.Sprintf("%s/%s.log", conf.LogPath, conf.AppName)
	fileWriter, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		j, _ := json.Marshal(conf)
		l.Fatal(context.TODO(), fmt.Sprintf("invalid log conf: %s, create log file %s failed: %v", j, logFile, err))
		return nil, err
	}

	l.WithWriter(fileWriter)
	return l, nil
}

func (l *Logger) WithWriter(wr io.Writer)            { l.wr = wr }
func (l *Logger) AddWriter(wrs ...io.Writer)         { l.wr = io.MultiWriter(l.wr, io.MultiWriter(wrs...)) }
func (l *Logger) WithReportHandler(rh ReportHandler) { l.reportFunc = rh }

////////////////////////////////////////////////////////////////////////////////
// 辅助方法

func (l *Logger) traceID(ctx context.Context) string {
	traceID := ctx.Value("trace_id")
	if traceID == nil {
		return ""
	}
	return fmt.Sprint(traceID)
}

func (l *Logger) formatLogMessage(ctx context.Context, withFileNum bool, prefix string, message string, vals ...interface{}) []byte {
	traceID := l.traceID(ctx)

	filenum := ""
	if withFileNum {
		if filenum = l.fileWithLineNum(); filenum != "" {
			filenum = "[" + filenum + "] "
		}
	}

	// [日志等级] [当前毫秒时间戳] [trace_id] [file:num] message append_message...
	msg := fmt.Sprintf("[%v] [%v] [%v] %v%v", prefix, time.Now().UnixNano()/1e6, traceID, filenum, message)
	for _, val := range vals {
		switch val := val.(type) {
		case string:
			msg = fmt.Sprintf("%v %s", msg, val)
		default:
			j, _ := json.Marshal(val)
			msg = fmt.Sprintf("%v %s", msg, j)
		}
	}

	return []byte(msg + "\n")
}

func (l *Logger) fileWithLineNum() string {
	substr := "/"
	if l.conf.AppName != "" {
		substr = "/" + l.conf.AppName + "/"
	}

	for i := 0; i < 17; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 跳过本文件
		if strings.HasSuffix(file, logSourceFile) {
			continue
		}

		if i := strings.LastIndex(file, substr); i > 0 {
			return file[i+1:] + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}

////////////////////////////////////////////////////////////////////////////////
// Write 实现 io.Writer 接口
// type Writer interface {
//     Write(p []byte) (n int, err error)
// }

// Write 对io.Writer接口的实现
func (l *Logger) Write(p []byte) (n int, err error) { return l.wr.Write(p) }

////////////////////////////////////////////////////////////////////////////////
// 接入 go-redis logger

// NewRedisLogger 新建go-redis日志实例
func NewRedisLogger(l *Logger) *log.Logger {
	switch l.conf.Level {
	case LOG_ALL, LOG_TRACE, LOG_DEBUG:
		// prefix ref: l.formatLogMessage(...)
		return log.New(l, "[info] [] [] go-redis ", log.LstdFlags|log.Lshortfile)
	default:
		// prefix ref: l.formatLogMessage(...)
		return log.New(l, "[info] [] [] go-redis ", log.LstdFlags)
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

// LogMode 不建议使用
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface { return l }

// Info -
func (l *Logger) Info(ctx context.Context, message string, vals ...interface{}) {
	if LOG_INFO > l.conf.Level {
		return
	}

	l.Write(l.formatLogMessage(ctx, l.conf.HasFileNum, "info", message, vals...))
}

// Warn -
func (l *Logger) Warn(ctx context.Context, message string, vals ...interface{}) {
	if LOG_WARN > l.conf.Level {
		return
	}

	l.Write(l.formatLogMessage(ctx, l.conf.HasFileNum, "warn", message, vals...))
}

// Error -
func (l *Logger) Error(ctx context.Context, message string, vals ...interface{}) {
	if LOG_ERROR > l.conf.Level {
		return
	}

	l.Write(l.formatLogMessage(ctx, l.conf.HasFileNum, "error", message, vals...))
}

// Trace -
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if LOG_TRACE > l.conf.Level {
		return
	}

	elapsed := float64(time.Since(begin).Nanoseconds()) / 1e6
	sql, rows := fc()

	format := ""
	if l.conf.HasColor {
		format = helper.ColorYellowBold + "[%vms] " + helper.ColorBlueBold + "[rows:%v] " + helper.ColorMagentaBold + "%v" + helper.ColorReset + ": %v"
	} else {
		format = "[%vms] [rows:%v] %v: %v"
	}
	msg := fmt.Sprintf(format, elapsed, rows, err, sql)

	l.Write(l.formatLogMessage(ctx, l.conf.HasFileNum, "trace", msg))
}

////////////////////////////////////////////////////////////////////////////////
// 补充实现其他方法

// Fatal -
func (l *Logger) Fatal(ctx context.Context, message string, vals ...interface{}) {
	if LOG_FATAL > l.conf.Level {
		return
	}

	logmsg := l.formatLogMessage(ctx, true, "fatal", message, vals...)

	// 写入日志
	l.Write(logmsg)

	// 触发告警
	if l.reportFunc != nil {
		arr := strings.Split(string(logmsg), " ")[3:]
		l.reportFunc(ctx, arr[0], strings.Join(arr[1:], " "))
	}
}

// Critical 同Fatal
func (l *Logger) Critical(ctx context.Context, message string, vals ...interface{}) {
	l.Fatal(ctx, message, vals...)
}

// Debug -
func (l *Logger) Debug(ctx context.Context, message string, vals ...interface{}) {
	if LOG_DEBUG > l.conf.Level {
		return
	}

	l.Write(l.formatLogMessage(ctx, l.conf.HasFileNum, "debug", message, vals...))
}

// Audit 审计
func (l *Logger) Audit(ctx context.Context, message string, vals ...interface{}) {
	if LOG_INFO > l.conf.Level {
		return
	}

	l.Write(l.formatLogMessage(ctx, l.conf.HasFileNum, "audit", message, vals...))
}
