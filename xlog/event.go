package xlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xnet"
)

var eventPool *sync.Pool
var appName string = filepath.Base(os.Args[0])
var localip string
var curPid int

func init() {
	eventPool = &sync.Pool{
		New: func() interface{} {
			return &Event{}
		},
	}
	localip = xnet.GetLocalIP()
	curPid = os.Getpid()
}

//Event 日志信息
type Event struct {
	Name    string
	Level   Level
	LogTime time.Time
	Session string
	Content string
	Output  string
	Tags    map[string]string
}

//NewEvent 构建日志事件
func NewEvent(name string, level Level, session string, content string, tags map[string]string) *Event {
	e := eventPool.Get().(*Event)
	e.LogTime = time.Now()
	e.Level = level
	e.Name = name
	e.Session = session
	e.Content = content
	e.Tags = tags
	return e
}

//Format 获取转换后的日志事件
func (e *Event) Format(layout *Layout) *Event {
	e.Output = e.Transform(layout.Layout, layout.IsJsonLayout)
	return e
}

func (e *Event) Transform(template string, isJson bool) string {
	word, _ := regexp.Compile(`%\w+`)

	//@变量, 将数据放入params中
	return word.ReplaceAllStringFunc(template, func(s string) string {
		key := s[1:]
		switch key {
		case "app":
			return appName
		case "nm":
			return e.Name
		case "session":
			return e.Session
		case "date":
			return e.LogTime.Format("20060102")
		case "datetime":
			return e.LogTime.Format("2006-01-02 15:04:05.000000")
		case "yy":
			return e.LogTime.Format("06")
		case "mm":
			return e.LogTime.Format("01")
		case "dd":
			return e.LogTime.Format("02")
		case "hh":
			return e.LogTime.Format("15")
		case "mi":
			return e.LogTime.Format("04")
		case "ss":
			return e.LogTime.Format("05")
		case "ms":
			return fmt.Sprintf("%06d", e.LogTime.Nanosecond()/1e3)
		case "level":
			return e.Level.FullName()
		case "l":
			return e.Level.Name()
		case "pid":
			return fmt.Sprintf("%d", curPid)
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
		case "ip":
			return localip
		case "tags":
			bytes, _ := json.Marshal(e.Tags)
			return bytesconv.BytesToString(bytes)
		default:
			return ""
		}
	})
}

//Close 关闭回收日志
func (e *Event) Close() {
	eventPool.Put(e)
}
