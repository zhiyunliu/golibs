package xtypes

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/zhiyunliu/golibs/xtransform"
)

type XMap map[string]interface{}

//Keys 从对象中获取数据值，如果不是字符串则返回空
func (m XMap) Keys() []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	return keys
}

//Merge 合并
func (m XMap) Merge(r XMap) {
	for k, v := range r {
		m[k] = v
	}
}

//Get 获取指定元素的值
func (m XMap) Get(name string) (interface{}, bool) {
	v, ok := m[name]
	return v, ok
}

//Get 获取指定元素的Bool值
func (m XMap) GetBool(name string) bool {
	v, ok := m[name]
	if ok {
		tmp, err := strconv.ParseBool(fmt.Sprint(v))
		if err != nil {
			return false
		}
		return tmp
	}
	return false
}

//Get 获取指定元素的值
func (m XMap) GetString(name string) string {
	v, ok := m[name]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%+v", v)
}

//Scan 以json 标签进行序列化
func (m XMap) Scan(obj interface{}) error {
	bytes, _ := json.Marshal(m)
	return json.Unmarshal(bytes, obj)
}

func (m XMap) Len() int {
	return len(m)
}

func (m XMap) IsEmpty() bool {
	return len(m) == 0
}

func (m XMap) SMap() SMap {
	sm := map[string]string{}
	for k, v := range m {
		sm[k] = fmt.Sprintf("%+v", v)
	}
	return sm
}

func (m XMap) Translate(tpl string) string {
	return xtransform.Translate(tpl, m)
}

func (m XMap) MarshalBinary() (data []byte, err error) {
	tmp := map[string]interface{}(m)
	return json.Marshal(tmp)
}

func (m XMap) GetInt(key string) (int, error) {
	tmp, ok := m[key]
	if !ok {
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

func (m XMap) GetInt64(key string) (int64, error) {
	tmp, ok := m[key]
	if !ok {
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

func (m XMap) GetFloat64(key string) (float64, error) {
	tmp, ok := m[key]
	if !ok {
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
