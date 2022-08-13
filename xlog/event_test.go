package xlog

import (
	"encoding/json"
	"testing"
	"text/template"
	"time"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/session"
)

func TestEvent_Transform(t *testing.T) {
	tests := []struct {
		name     string
		e        *Event
		template string
		isJson   bool
		want     string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Transform(tt.template, tt.isJson); got != tt.want {
				t.Errorf("Event.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_Event_Transform(b *testing.B) {
	evt := &Event{
		Level:   LevelInfo,
		Session: session.Create(),
		Idx:     1,
		Content: "Benchmark_Event_Transform",
	}
	template := "[%time][%l][%session][%idx] %content"
	isJson := false
	for i := 0; i < b.N; i++ {
		evt.Transform(template, isJson)
	}
}

func Benchmark_Event_Transform2(b *testing.B) {
	evt := &Event{
		LogTime: time.Now(),
		Level:   LevelInfo,
		Session: session.Create(),
		Idx:     1,
		Content: "i am content.",
	}
	content := "[{{XTime .LogTime}}][{{.Level}}][{{.Session}}][{{.Idx}}] {{content .Content false}}\n"
	funcMap := template.FuncMap{}
	funcMap["content"] = func(content string, isJson bool) string {
		if isJson {
			bytes, _ := json.Marshal(content)
			return bytesconv.BytesToString(bytes)
		}
		return content
	}
	funcMap["XTime"] = func(logTime time.Time) string {
		return logTime.Format("15:04:05.000000")
	}
	// funcMap["Time"] = func(logTime time.Time) string {
	// 	return logTime.Format("15:04:05.000000")
	// }
	tpl, err := template.New("xlog").Funcs(funcMap).Parse(content)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		tpl.Execute(&noneWriter{}, evt)
	}
}

type noneWriter struct {
}

func (w *noneWriter) Write(p []byte) (n int, err error) {
	//fmt.Fprint(os.Stdout, string(p))
	return
}
