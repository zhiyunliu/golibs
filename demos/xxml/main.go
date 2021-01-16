package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/lib4go/file"
	"github.com/zhiyunliu/golibs/types/xxml"
)

type Val struct {
	Name  string
	Value string
}

type DemoStruct struct {
	F1  string
	F2  string `xml:"f2"`
	F3  int
	F4  map[string]string
	F5  map[string]interface{}
	F6  map[int]string
	F7  XMap
	F8  *XMap
	F9  []string
	F10 []Val
	F11 []*Val
	F12 IntArray
	F13 *time.Time
}
type XMap map[string]interface{}

func main() {
	m := map[string]interface{}{
		/*		"a": 1,
				"b": 3.2,
				"c": DemoStruct{
					F1: "f1",
					F2: "f2",
					F3: 3,
					F4: map[string]string{"f4": "f4"},
					F5: map[string]interface{}{"f5": 5},
					F6: map[int]string{1: "1", 2: "f6"},
				},
				"d": &DemoStruct{
					F1: "df1",
					F7: XMap{
						"df1": "1",
					},
					F8: &XMap{
						"df2": 2,
					},
				},
				"e": &DemoStruct{
					F9:  []string{"a", "b", "c"},
					F10: []Val{{Name: "name", Value: "val"}, {Name: "n2"}},
					F11: []*Val{{Name: "name", Value: "val"}, {Name: "n2"}},
				},
		*/
		"f": IntArray{{K: "1", V: "2"}, {K: "3", V: "4"}},
		"g": time.Now(),
	}
	val, err := xxml.Marshal(m, xxml.WithUncompress(), xxml.WithCustomType(&IntArray{}))
	f, err := file.CreateFile("xx")
	fmt.Fprintln(f, val)
	fmt.Fprintln(f, err)

}
