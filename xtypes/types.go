package xtypes

import (
	"encoding/json"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xreflect"
)

var (
	// Deprecated: As of Go v0.2.0, this function simply calls [xreflect.GetBool].
	GetBool = xreflect.GetBool
	// Deprecated: As of Go v0.2.0, this function simply calls [xreflect.GetString].
	GetString = xreflect.GetString
	// Deprecated: As of Go v0.2.0, this function simply calls [xreflect.GetInt].
	GetInt = xreflect.GetInt
	// Deprecated: As of Go v0.2.0, this function simply calls [xreflect.GetInt64].
	GetInt64 = xreflect.GetInt64
	// Deprecated: As of Go v0.2.0, this function simply calls [xreflect.GetUint64].
	GetUint64 = xreflect.GetUint64
	// Deprecated: As of Go v0.2.0, this function simply calls [xreflect.GetFloat64].
	GetFloat64 = xreflect.GetFloat64
)

func mapscan(obj any, m any) error {

	if obj == nil {
		return nil
	}

	switch v := obj.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal(bytesconv.StringToBytes(v), m)

	case *[]byte:
		if v == nil {
			return nil
		}
		return json.Unmarshal(*v, m)
	case *string:
		if v == nil {
			return nil
		}
		return json.Unmarshal(bytesconv.StringToBytes(*v), m)
	}
	return nil
}
