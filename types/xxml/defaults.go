package xxml

import (
	"fmt"
)

var defaults = &Options{
	RootElement:    "xml",
	DefaultElement: "item",
	TimeFormat:     "2006-01-02 15:04:05",
	NumberType:     "decimal",
}

func errorUnsupportType(typename string) error {
	return fmt.Errorf("unsupport types:%s", typename)
}
