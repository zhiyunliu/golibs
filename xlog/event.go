package xlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/zhiyunliu/golibs/bytesconv"
)

var (
	curPid string

	eventPool *sync.Pool
	appName   string = filepath.Base(os.Args[0])
	word, _          = regexp.Compile(`%\w+`)
)

func init() {
	curPid = strconv.FormatInt(int64(os.Getpid()), 10)

	eventPool = &sync.Pool{
		New: func() interface{} {
			return &Event{}
		},
	}
}

// Event 日志信息
type Event struct {
	Name    string
	Level   Level
	Idx     int32
	SrvType string
	LogTime time.Time
	Session string
	Content string
	Output  string
	Tags    map[string]string
}

// NewEvent 构建日志事件
func NewEvent(name string, level Level, session string, srvType string, content string, tags map[string]string) *Event {
	e := eventPool.Get().(*Event)
	e.LogTime = time.Now()
	e.Level = level
	e.Name = name
	e.Session = session
	e.Content = content
	e.SrvType = srvType
	e.Tags = tags
	return e
}

// Format 获取转换后的日志事件
func (e *Event) Format(layout *Layout) *Event {
	e.Output = e.Transform(layout.Content, layout.isJson)
	return e
}

func (e *Event) Transform(template string, isJson bool) string {
	//@变量, 将数据放入params中
	return word.ReplaceAllStringFunc(template, func(s string) string {
		key := s[1:]
		switch key {
		case "app":
			return appName
		case "pid":
			return curPid
		case "nm":
			return e.Name
		case "srvtype":
			return e.SrvType
		case "session":
			return e.Session
		case "date":
			return e.LogTime.Format("20060102")
		case "ndate":
			return e.LogTime.Format("2006-01-02")
		case "time":
			return e.LogTime.Format("15:04:05.000000")
		case "datetime":
			return e.LogTime.Format("2006-01-02 15:04:05.000000")
		case "yy":
			return strconv.FormatInt(int64(e.LogTime.Year()), 10)
		case "mm":
			return strconv.FormatInt(int64(e.LogTime.Month()), 10)
		case "dd":
			return strconv.FormatInt(int64(e.LogTime.Day()), 10)
		case "hh":
			return strconv.FormatInt(int64(e.LogTime.Hour()), 10)
		case "mi":
			return strconv.FormatInt(int64(e.LogTime.Minute()), 10)
		case "ss":
			return strconv.FormatInt(int64(e.LogTime.Second()), 10)
		case "ms":
			return fmt.Sprintf("%06d", e.LogTime.Nanosecond()/1e3)
		case "level":
			return e.Level.FullName()
		case "l":
			return e.Level.Name()
		case "idx":
			return strconv.FormatInt(int64(e.Idx), 10)
		case "n":
			return "\n"
		case "content":
			if isJson {
				buff, err := json.Marshal(e.Content)
				if err != nil {
					return e.Content
				}
				if len(buff) > 2 {
					return bytesconv.BytesToString(buff[1 : len(buff)-1])
				}
				return bytesconv.BytesToString(buff)
			}
			return e.Content
		default:
			return Transform(key, e, isJson)
		}
	})
}

// Close 关闭回收日志
func (e *Event) Close() {
	eventPool.Put(e)
}
