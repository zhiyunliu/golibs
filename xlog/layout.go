package xlog

import (
	"encoding/json"
	"fmt"

	"github.com/zhiyunliu/golibs/xfile"
)

//Layout 输出器
type Layout struct {
	Type         string `json:"type"`
	LevelName    string `json:"level" valid:"in(off|info|warn|error|fatal|debug|all)"`
	Path         string `json:"path,omitempty"`
	Layout       string `json:"layout"`
	Level        Level  `json:"-"`
	IsJsonLayout bool   `json:"-"`
}

func (l *Layout) Init() {
	l.IsJsonLayout = json.Valid([]byte(l.Layout))
	l.Level = TransLevel(l.LevelName)
}

type layoutSetting struct {
	Status  bool      `json:""`
	Layouts []*Layout `json:"layouts" toml:"layouts"`
}

func newDefLayouts() *layoutSetting {
	setting := &layoutSetting{Layouts: make([]*Layout, 0, 2)}
	defaultLayout := "[%datetime.%ms][%l][%session][%data] %content%n"

	fileLayout := &Layout{Type: File, LevelName: LevelAll.Name()}
	fileLayout.Path, _ = xfile.GetAbs("../logs/%app/%date/%hh.log")
	fileLayout.Layout = defaultLayout
	fileLayout.Init()
	setting.Layouts = append(setting.Layouts, fileLayout)

	stdLayout := &Layout{Type: Stdout, LevelName: LevelAll.Name()}
	stdLayout.Layout = defaultLayout
	stdLayout.Init()
	setting.Layouts = append(setting.Layouts, stdLayout)

	return setting
}

//Encode 将当前配置内容保存到文件中
func Encode(path string) error {
	f, err := xfile.CreateFile(path)
	if err != nil {
		return fmt.Errorf("无法创建文件:%s %w", path, err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	err = encoder.Encode(newDefLayouts())
	if err != nil {
		return err
	}
	return nil
}

//Decode 从配置文件中读取配置信息
func Decode(path string) (*layoutSetting, error) {
	l := &layoutSetting{}
	f, err := xfile.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取文件:%s %w", path, err)
	}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(l); err != nil {
		return nil, err
	}
	return l, nil
}