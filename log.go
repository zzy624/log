package log

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var _log *log.Entry

type Fields log.Fields

func GetLog(Topic string) *log.Entry {
	zlog := log.New()
	_log := log.NewEntry(zlog)
	zlog.SetFormatter(&log.TextFormatter{
		QuoteEmptyFields:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		ForceColors:            true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, file, line, _ := runtime.Caller(8)
			files := strings.Split(file, "/")
			if len(files) > 2 {
				file = strings.Join(files[len(files)-2:], "/")
			} else {
				file = strings.Join(files, "/")
			}
			return file, strconv.Itoa(line)
		},
		FieldMap: log.FieldMap{
			FieldKeyTime:  "timestamp",
			FieldKeyLevel: "level",
			FieldKeyMsg:   "message",
			FieldKeyFunc:  "caller",
		},
	})
	_log = zlog.WithFields(log.Fields{FieldkeyTopic: Topic})
	zlog.AddHook(newLfsHook(zlog, Topic, log.DebugLevel, 7))

	return _log
}

//对外提供统一接口，可自定义替换
//默认使用dlog
type Logger interface {
	Debug(kv ...interface{})
	Info(kv ...interface{})
	Warn(kv ...interface{})
	Error(kv ...interface{})
	Panic(kv ...interface{})
	WithFields(fields log.Fields) *log.Entry
	//Close() error
	//DebugLog(b bool)
}

func SetTopic(Topic string) {
	_log = GetLog(Topic)
}

func SetLogger(l *log.Entry) {
	_log = l
}
func GetLogger() Logger {
	if _log == nil {
		topic := "default_server"
		SetLogger(GetLog(topic))
	}
	return _log
}

func Debug(msg interface{}, kv ...interface{}) {
	if len(kv) == 0 {
		GetLogger().Debug(msg)
	} else {
		_fields := make(log.Fields, 0)
		if len(kv)%2 != 0 {
			kv = append(kv, "unknown")
		}
		for i := 0; i < len(kv); i += 2 {
			_fields[fmt.Sprintf("%v", kv[i])] = kv[i+1]
		}
		WithFields(_fields).Debug(msg)
	}

}
func Info(msg interface{}, kv ...interface{}) {
	if len(kv) == 0 {
		GetLogger().Info(msg)
	} else {
		_fields := make(log.Fields, 0)
		if len(kv)%2 != 0 {
			kv = append(kv, "unknown")
		}
		for i := 0; i < len(kv); i += 2 {
			_fields[fmt.Sprintf("%s", kv[i])] = kv[i+1]
		}
		WithFields(_fields).Info(msg)
	}

}
func Warn(msg interface{}, kv ...interface{}) {
	if len(kv) == 0 {
		GetLogger().Warn(msg)
	} else {
		_fields := make(log.Fields, 0)
		if len(kv)%2 != 0 {
			kv = append(kv, "unknown")
		}
		for i := 0; i < len(kv); i += 2 {
			_fields[fmt.Sprintf("%s", kv[i])] = kv[i+1]
		}
		WithFields(_fields).Warn(msg)
	}
}
func Error(msg interface{}, kv ...interface{}) {
	if len(kv) == 0 {
		GetLogger().Error(msg)
	} else {
		_fields := make(log.Fields, 0)
		if len(kv)%2 != 0 {
			kv = append(kv, "unknown")
		}
		for i := 0; i < len(kv); i += 2 {
			_fields[fmt.Sprintf("%s", kv[i])] = kv[i+1]
		}
		WithFields(_fields).Error(msg)
	}
}
func Panic(msg interface{}, kv ...interface{}) {
	if len(kv) == 0 {
		GetLogger().Panic(msg)
	} else {
		_fields := make(log.Fields, 0)
		if len(kv)%2 != 0 {
			kv = append(kv, "unknown")
		}
		for i := 0; i < len(kv); i += 2 {
			_fields[fmt.Sprintf("%s", kv[i])] = kv[i+1]
		}
		WithFields(_fields).Panic(msg)
	}
}

func WithFields(_fields log.Fields) *log.Entry {
	return GetLogger().WithFields(_fields)
}

func newLfsHook(zlog *log.Logger, logName string, logLevel log.Level, maxRemainCnt uint) log.Hook {
	writer, err := rotatelogs.New(
		logName+".%Y%m%d%H%M",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour*24),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}

	zlog.SetLevel(logLevel)
	zlog.SetReportCaller(true)

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &JSONFormatter{
		DataKey:         "data",
		//PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(_func, file string) (string, string) {
			files := strings.Split(_func, "/")
			if len(files) > 2 {
				_func = strings.Join(files[len(files)-2:], "/")
			} else {
				_func = strings.Join(files, "/")
			}
			files2 := strings.Split(file, "/")
			if len(files2) > 2 {
				file = strings.Join(files2[len(files2)-2:], "/")
			} else {
				file = strings.Join(files2, "/")
			}
			return _func, file
		},
		FieldMap: FieldMap{
			FieldKeyTime:  "timestamp",
			FieldKeyLevel: "level",
			FieldKeyMsg:   "msg",
			FieldKeyFunc:  "caller",
		}})

	return lfsHook
}
