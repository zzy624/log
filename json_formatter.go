package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
)

type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

type logStore struct {
	Topic     interface{} `json:"topic"`
	Level     interface{} `json:"level"`
	Func      interface{} `json:"func"`
	File      interface{} `json:"file"`
	Line      interface{} `json:"line"`
	Msg       interface{} `json:"msg"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp interface{} `json:"timestamp"`
}

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime:  "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg:   "@message",
	// 		 FieldKeyFunc:  "@caller",
	//    },
	// }
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(string,string) (_func,file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	data := make(log.Fields, len(entry.Data)+4)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	var topic interface{}
	//var _data log.Fields

	if len(data) > 0 {
		topic = data[FieldkeyTopic]
		delete(data, FieldkeyTopic)
	}

	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	var l logStore
	l.Topic = topic
	l.Level = entry.Level.String()
	l.Msg = entry.Message
	if len(data) > 0 {
		l.Data = data
	}
	l.Timestamp = entry.Time.Format(timestampFormat)

	pc, file, line, _ := runtime.Caller(9)
	var _func string
	_f := runtime.FuncForPC(pc)
	_func = _f.Name()

	if f.CallerPrettyfier != nil {
		_func,file = f.CallerPrettyfier(_func,file)
	}
	if file != "" {
		l.File = file
	}
	if line != 0 {
		l.Line = line
	}
	if _func != "" {
		l.Func = _func
	}

	encoder := json.NewEncoder(b)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(l); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}

	return b.Bytes(), nil
}
