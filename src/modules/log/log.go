package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	baseLog "github.com/weisd/log"
)

type LogService struct {
	Enable  bool
	Mode    string
	Level   string
	BuffLen int64

	// model file
	FileName     string // 文件名
	LogRotate    bool   // 分割文件
	MaxLines     int    // 最大行数
	MaxSizeShift int    // 最大文件大小 1 << MaxSizeShift
	DailyRotate  bool   // 每天分割文件
	MaxDays      int    // 分割文件保留天数

	// model conn
	ReConnetOnMsg bool
	ReConnet      bool
	Protocol      string // tcp udp unix
	Addr          string

	// model smtp
	User      string
	Passwd    string
	Host      string
	Receivers []string // default "[]"
	Subject   string

	// database
	Driver string
	Conn   string
}

// 所有的log
var LogsMap map[string]loggerMap

var logLevels = map[string]string{
	"trace":    "0",
	"debug":    "1",
	"info":     "2",
	"warn":     "3",
	"error":    "4",
	"critical": "5",
}

// func init() {
// 	baseLog.NewLogger(0, "console", `{"level":0}`)
// }

func InitLogService(logsMaps map[string]map[string]LogService) {
	LogsMap = make(map[string]loggerMap)

	for name, m := range logsMaps {
		loggers := make(loggerMap)
		for k, v := range m {
			if !v.Enable {
				continue
			}
			str := ""

			level := logLevels[strings.ToLower(v.Level)]
			if len(level) == 0 {
				continue
			}

			switch v.Mode {
			case "console":
				str = fmt.Sprintf(`{"level":%s}`, level)
			case "file":
				filename := "./log/app.log"
				logRotate := true
				maxLine := 10000
				maxSizeShift := 28
				dailyRotate := true
				maxDays := 7
				if len(v.FileName) > 0 {
					filename = v.FileName
				}
				// 创建目录
				if len(filename) > 0 {
					logpath, _ := filepath.Abs(filename)
					os.MkdirAll(path.Dir(logpath), os.ModePerm)

				}

				if !v.LogRotate {
					logRotate = false
				}
				if v.MaxLines > 0 {
					maxLine = v.MaxLines
				}
				if v.MaxSizeShift > 0 {
					maxSizeShift = v.MaxSizeShift
				}
				if !v.DailyRotate {
					dailyRotate = false
				}
				if v.MaxDays > 0 {
					maxDays = v.MaxDays
				}
				str = fmt.Sprintf(
					`{"level":%s,"filename":"%s","rotate":%v,"maxlines":%d,"maxsize":%d,"daily":%v,"maxdays":%d}`,
					level,
					filename,
					logRotate,
					maxLine,
					1<<uint(maxSizeShift),
					dailyRotate, maxDays,
				)
			case "conn":
				str = fmt.Sprintf(`{"level":%s,"reconnectOnMsg":%v,"reconnect":%v,"net":"%s","addr":"%s"}`, level,
					v.ReConnetOnMsg,
					v.ReConnet,
					v.Protocol,
					v.Addr,
				)
			case "smtp":
				tos, err := json.Marshal(v.Receivers)
				if err != nil {

					baseLog.Error(4, "json.Marshal(v.Receivers) err %v", err)
					continue
				}

				str = fmt.Sprintf(`{"level":%s,"username":"%s","password":"%s","host":"%s","sendTos":%s,"subject":"%s"}`, level,
					v.User,
					v.Passwd,
					v.Host,
					tos,
					v.Subject)

			default:
				baseLog.Warn("unkown log mode %s", v.Mode)
				continue
			}

			var bufflen int64 = 10000
			if v.BuffLen > 0 {
				bufflen = v.BuffLen
			}

			loggers[k] = baseLog.NewCustomLogger(bufflen, v.Mode, str)
		}

		LogsMap[name] = loggers
	}

	// baseLog.Debug("init logs done %v", LogsMap)
}

///////////// logger //////////////

type loggerMap map[string]*baseLog.Logger

func (loggers loggerMap) Trace(format string, v ...interface{}) {

	for _, logger := range loggers {
		logger.Trace(format, v...)
	}
}

func (loggers loggerMap) Debug(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Debug(format, v...)
	}
}

func (loggers loggerMap) Info(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Info(format, v...)
	}
}

func (loggers loggerMap) Warn(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Warn(format, v...)
	}
}

func (loggers loggerMap) Error(skip int, format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Error(skip, format, v...)
	}
}

func (loggers loggerMap) Critical(skip int, format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Critical(skip, format, v...)
	}
}

func (loggers loggerMap) Fatal(skip int, format string, v ...interface{}) {
	loggers.Error(skip, format, v...)
	for _, l := range loggers {
		l.Close()
	}
	os.Exit(1)
}

func Get(name string) loggerMap {
	l, ok := LogsMap[name]
	if !ok {
		baseLog.Fatal(4, "Unknown log %s", name)
		return nil
	}

	return l
}

func Trace(format string, v ...interface{}) {

	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Trace(format, v...)
}

func Debug(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Error(4, format, v...)
}

func Critical(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Critical(4, format, v...)
}

func Fatal(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Fatal(4, format, v...)
}

func Close(name ...string) {
	k := "default"
	if len(name) > 0 {
		k = name[0]
	}

	for _, l := range LogsMap[k] {
		l.Close()
	}

}
