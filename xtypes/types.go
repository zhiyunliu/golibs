package xtypes

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/zhiyunliu/golibs/bytesconv"
)

func GetBool(tmp interface{}) bool {
	if tmp == nil {
		return false
	}
	tmpB, err := strconv.ParseBool(fmt.Sprint(tmp))
	if err != nil {
		return false
	}
	return tmpB
}

func GetString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%+v", v)
}

func GetInt(tmp interface{}) (int, error) {

	if tmp == nil {
		return 0, nil
	}
	switch val := tmp.(type) {
	case *int:
		return *val, nil
	case *int8:
		return int(*val), nil
	case *int16:
		return int(*val), nil
	case *int32:
		return int(*val), nil
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		if math.MinInt <= val && val <= math.MaxInt {
			return int(val), nil
		}
		return 0, fmt.Errorf("数据越界:int64=>int,%d", val)
	case string:
		return strToint(val)
	default:
		if strer, ok := tmp.(fmt.Stringer); ok {
			return strToint(strer.String())
		}
	}

	return 0, newNotSupportErr(tmp)
}

func GetInt64(tmp interface{}) (int64, error) {
	if tmp == nil {
		return 0, nil
	}
	switch val := tmp.(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case *int:
		return int64(*val), nil
	case *int8:
		return int64(*val), nil
	case *int16:
		return int64(*val), nil
	case *int32:
		return int64(*val), nil
	case *int64:
		return *val, nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		if strer, ok := tmp.(fmt.Stringer); ok {
			return strconv.ParseInt(strer.String(), 10, 64)
		}
	}

	return 0, newNotSupportErr(tmp)
}

func GetUint64(tmp interface{}) (uint64, error) {
	if tmp == nil {
		return 0, nil
	}
	switch val := tmp.(type) {
	case uint:
		return uint64(val), nil
	case uint8:
		return uint64(val), nil
	case uint16:
		return uint64(val), nil
	case uint32:
		return uint64(val), nil
	case uint64:
		return val, nil

	case *uint:
		return uint64(*val), nil
	case *uint8:
		return uint64(*val), nil
	case *uint16:
		return uint64(*val), nil
	case *uint32:
		return uint64(*val), nil
	case *uint64:
		return *val, nil

	case string:
		return strconv.ParseUint(val, 10, 64)
	default:
		if strer, ok := tmp.(fmt.Stringer); ok {
			return strconv.ParseUint(strer.String(), 10, 64)
		}
	}

	return 0, newNotSupportErr(tmp)
}

func GetFloat64(tmp interface{}) (float64, error) {
	if tmp == nil {
		return 0, nil
	}
	switch val := tmp.(type) {
	case float32:
		return float64(val), nil
	case float64:
		return float64(val), nil
	case *float32:
		return float64(*val), nil
	case *float64:
		return float64(*val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		if strer, ok := tmp.(fmt.Stringer); ok {
			return strconv.ParseFloat(strer.String(), 64)
		}
	}

	return 0, newNotSupportErr(tmp)
}

func strToint(str string) (int, error) {
	var t64 int64
	t64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	if math.MinInt <= t64 && t64 <= math.MaxInt {
		return int(t64), nil
	}
	return 0, fmt.Errorf("数据越界:int64=>int,%d", t64)
}

func newNotSupportErr(tmp any) error {
	return fmt.Errorf("不支持的数据类型:%s", reflect.TypeOf(tmp).Name())
}

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
