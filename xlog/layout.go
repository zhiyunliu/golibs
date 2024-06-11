package xlog

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zhiyunliu/golibs/xfile"
)

// Layout 输出器
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

type LayoutSetting struct {
	Enable bool               `json:"enable"`
	Layout map[string]*Layout `json:"layout"`
}

// Encode 将当前配置内容保存到文件中
func Encode(path string, setting *LayoutSetting) error {
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

// Decode 从配置文件中读取配置信息
func Decode(path string) (*LayoutSetting, error) {
	setting := &LayoutSetting{
		Enable: true,
		Layout: map[string]*Layout{},
	}
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败:%s %w", path, err)
	}
	err = json.Unmarshal(fileBytes, setting)
	return setting, err
}

func loadLayout(paths ...string) (setting *LayoutSetting, err error) {
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
		setting = &LayoutSetting{Enable: true, Layout: map[string]*Layout{}}
		_appenderCache.Range(func(key, value interface{}) bool {
			setting.Layout[fmt.Sprintf("%s", key)] = DefaultParam.Layout
			return true
		})
	} else {
		setting, err = Decode(path)
		if err != nil {
			err = fmt.Errorf("读取配置文件失败 %v", err)
			return
		}
	}
	return
}
