package xlog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/zhiyunliu/golibs/xfile"
)

//Layout 输出器
type Layout struct {
	//Type      string `json:"type"`
	LevelName string `json:"level" valid:"in(off|info|warn|error|panic|fatal|debug|all)"`
	Path      string `json:"path,omitempty"`
	Content   string `json:"content"`
	Level     Level  `json:"-"`
	isJson    bool   `json:"-"`
}

func (l *Layout) Init() {
	l.isJson = json.Valid([]byte(l.Content))
	l.Level = TransLevel(l.LevelName)
}

type layoutSetting struct {
	Enable bool               `json:"enable"`
	Layout map[string]*Layout `json:"layout"`
}

//Encode 将当前配置内容保存到文件中
func Encode(path string, setting *layoutSetting) error {
	f, err := xfile.CreateFile(path)
	if err != nil {
		return fmt.Errorf("无法创建文件:%s %w", path, err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	err = encoder.Encode(setting)
	if err != nil {
		return err
	}
	return nil
}

//Decode 从配置文件中读取配置信息
func Decode(path string) (*layoutSetting, error) {
	l := &layoutSetting{}
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBytes, l)
	return l, err
}

func loadLayout(paths ...string) (setting *layoutSetting, err error) {
	if len(paths) <= 0 {
		return nil, fmt.Errorf("未设置日志文件配置路径")
	}

	var path string = paths[0]
	for _, tmp := range paths {
		if xfile.Exists(tmp) {
			path = tmp
			break
		}
	}

	if !xfile.Exists(path) {
		setting = &layoutSetting{Enable: true, Layout: map[string]*Layout{}}
		_appenderCache.Range(func(key, value interface{}) bool {
			builder := value.(AppenderBuilder)
			name := fmt.Sprintf("%s", key)
			setting.Layout[name] = builder.DefaultLayout()
			return true
		})

		err = Encode(path, setting)
		if err != nil {
			err = fmt.Errorf("创建日志配置文件失败 %v", err)
			return
		}
	} else {
		setting, err = Decode(path)
		if err != nil {
			err = fmt.Errorf("读取配置文件失败 %v", err)
			return
		}
	}

	_globalPause = !setting.Enable
	return
}
