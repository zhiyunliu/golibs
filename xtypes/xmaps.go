package xtypes

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xtypes/internal"
)

type XMaps []XMap

func (ms *XMaps) Append(i ...XMap) XMaps {
	*ms = append(*ms, i...)
	return *ms
}

func (ms XMaps) IsEmpty() bool {
	return len(ms) == 0
}

func (ms XMaps) Len() int {
	return len(ms)
}

func (ms XMaps) Get(idx int) XMap {
	if idx < 0 || len(ms) <= idx {
		return map[string]interface{}{}
	}
	return ms[idx]
}

// Deprecated: As of Go v0.2.0, this function simply calls [ScanTo].
func (ms XMaps) Scan(obj interface{}) error {
	return ms.ScanTo(obj)
}

func (ms XMaps) ScanTo(obj interface{}) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &internal.XMapScanError{Type: reflect.TypeOf(obj)}
	}

	sliceType := rv.Elem().Type()

	if !(sliceType.Kind() == reflect.Slice) {
		return &internal.XMapScanError{Message: "Not Slice"}
	}
	structType := sliceType.Elem()
	var subIsPtr bool = false
	if structType.Kind() == reflect.Ptr {
		subIsPtr = true
		structType = structType.Elem()
	}

	newv := reflect.MakeSlice(sliceType, len(ms), len(ms))
	for i := range ms {

		tmp := reflect.New(structType)
		err := ms[i].ScanTo(tmp.Interface())
		if err != nil {
			return err
		}
		itemV := tmp
		if !subIsPtr {
			itemV = itemV.Elem()
		}
		newv.Index(i).Set(itemV)
	}

	//reflect.Copy(rv.Elem(), newv)
	rv.Elem().Set(newv)
	return nil
}

func (m *XMaps) MapScan(obj interface{}) error {
	if obj == nil {
		return nil
	}
	*m = XMaps{}
	return mapscan(obj, m)
}

func (m XMaps) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// Value String
func (m XMaps) Value() (driver.Value, error) {
	bytes, err := m.MarshalBinary()
	return bytesconv.BytesToString(bytes), err
}
