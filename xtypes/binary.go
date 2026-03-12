package xtypes

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/zhiyunliu/golibs/bytesconv"
)

type Binary []byte

func (b *Binary) Scan(src any) error {
	if src == nil {
		return nil
	}
	*b = src.([]byte)
	return nil
}

func (b Binary) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte("null"), nil
	}
	return json.Marshal(bytesconv.BytesToString(b))
}

func (b Binary) Value() (driver.Value, error) {
	return b, nil
}
