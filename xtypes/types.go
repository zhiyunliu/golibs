package xtypes

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
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
	case int:
		return val, nil
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

	return 0, fmt.Errorf("不支持的数据类型:%s", reflect.TypeOf(tmp).Name())
}

func GetInt64(tmp interface{}) (int64, error) {
	if tmp == nil {
		return 0, nil
	}

	switch val := tmp.(type) {
	case int:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case string:
		return strToint64(val)
	default:
		if strer, ok := tmp.(fmt.Stringer); ok {
			return strToint64(strer.String())
		}
	}

	return 0, fmt.Errorf("不支持的数据类型:%s", reflect.TypeOf(tmp).Name())
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
	case string:
		return strTofloat64(val)
	default:
		if strer, ok := tmp.(fmt.Stringer); ok {
			return strTofloat64(strer.String())
		}
	}

	return 0, fmt.Errorf("不支持的数据类型:%s", reflect.TypeOf(tmp).Name())
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

func strToint64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
func strTofloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}
