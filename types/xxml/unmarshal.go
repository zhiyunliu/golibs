package xxml

import (
	"encoding/xml"
	"strings"
)

//Unmarshal 反序列化对象
func Unmarshal(val string, result interface{}, opts ...Option) error {

	decoder := xml.NewDecoder(strings.NewReader(val))

	decoder.Decode(result)

	return nil
}
