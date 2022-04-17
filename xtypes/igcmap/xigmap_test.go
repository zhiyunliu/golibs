package igcmap

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		orignal map[string]interface{}
		want    *IgcMap
	}{
		{name: "0.Nil初始化", orignal: nil, want: &IgcMap{data: map[string]interface{}{}}},
		{name: "1.无重复", orignal: map[string]interface{}{"a": 1, "b": "2"}, want: &IgcMap{data: map[string]interface{}{"a": 1, "b": "2"}}},
		{name: "2.有重复,按key排序取后者", orignal: map[string]interface{}{"a": 1, "A": "2"}, want: &IgcMap{data: map[string]interface{}{"a": 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.orignal); !reflect.DeepEqual(got.data, tt.want.data) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIgcMap_Set(t *testing.T) {
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name   string
		m      *IgcMap
		args   args
		keylen int
		want   bool
	}{
		{name: "1.不存在", m: New(map[string]interface{}{}), args: args{key: "a", val: 1}, want: false, keylen: 1},
		{name: "2.不存一致的key", m: New(map[string]interface{}{"a": 1}), args: args{key: "b", val: "2"}, want: false, keylen: 2},
		{name: "3.存在相同的key", m: New(map[string]interface{}{"a": 1}), args: args{key: "a", val: 2}, want: true, keylen: 1},
		{name: "4.存在大小写不一致的key", m: New(map[string]interface{}{"a": 1}), args: args{key: "A", val: "2"}, want: true, keylen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Set(tt.args.key, tt.args.val)
			if got != tt.want {
				t.Errorf("IgcMap.Set() = %v, want %v", got, tt.want)
			}
			if len(tt.m.keyMap) != tt.keylen {
				t.Errorf("IgcMap.Set() KeyLen %v", len(tt.m.keyMap))
			}
		})
	}
}

func TestIgcMap_Get(t *testing.T) {
	tests := []struct {
		name    string
		m       *IgcMap
		key     string
		wantVal interface{}
		wantOk  bool
	}{
		{name: "1.不存在的key", m: New(map[string]interface{}{}), key: "a", wantVal: nil, wantOk: false},
		{name: "2.存在大小写不一致的key", m: New(map[string]interface{}{"A": 1}), key: "a", wantVal: 1, wantOk: true},
		{name: "3.存在大小写一致的key", m: New(map[string]interface{}{"a": 1}), key: "a", wantVal: 1, wantOk: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := tt.m.Get(tt.key)
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("IgcMap.Get() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOk != tt.wantOk {
				t.Errorf("IgcMap.Get() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestIgcMap_MergeMap(t *testing.T) {
	tests := []struct {
		name   string
		m      *IgcMap
		args   map[string]interface{}
		keyLen int
	}{
		{name: "1.空合并", m: New(map[string]interface{}{}), args: map[string]interface{}{}, keyLen: 0},
		{name: "2.合并新值", m: New(map[string]interface{}{}), args: map[string]interface{}{"a": 1}, keyLen: 1},
		{name: "3.覆盖为新值", m: New(map[string]interface{}{"a": 1}), args: map[string]interface{}{"a": 2}, keyLen: 1},
		{name: "4.添加新值", m: New(map[string]interface{}{"a": 1}), args: map[string]interface{}{"b": 2}, keyLen: 2},
		{name: "5.大小写不一致", m: New(map[string]interface{}{"a": 1}), args: map[string]interface{}{"A": 10}, keyLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.MergeMap(tt.args)
			if len(tt.m.data) != tt.keyLen {
				t.Errorf("MergeMap gotLen = %v, want %v", len(tt.m.data), tt.keyLen)
			}
		})
	}
}

func TestIgcMap_MergeIgc(t *testing.T) {

	tests := []struct {
		name   string
		m      *IgcMap
		args   *IgcMap
		keyLen int
	}{
		{name: "1.空合并", m: New(map[string]interface{}{}), args: New(map[string]interface{}{}), keyLen: 0},
		{name: "2.合并新值", m: New(map[string]interface{}{}), args: New(map[string]interface{}{"a": 1}), keyLen: 1},
		{name: "3.覆盖为新值", m: New(map[string]interface{}{"a": 1}), args: New(map[string]interface{}{"a": 2}), keyLen: 1},
		{name: "4.添加新值", m: New(map[string]interface{}{"a": 1}), args: New(map[string]interface{}{"b": 2}), keyLen: 2},
		{name: "5.大小写不一致", m: New(map[string]interface{}{"a": 1}), args: New(map[string]interface{}{"A": 10}), keyLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.MergeIgc(tt.args)
			if len(tt.m.data) != tt.keyLen {
				t.Errorf("MergeMap gotLen = %v, want %v", len(tt.m.data), tt.keyLen)
			}
		})
	}
}

func TestIgcMap_Iter(t *testing.T) {
	tests := []struct {
		name     string
		m        *IgcMap
		callback func(key string, val interface{}) bool
	}{
		{name: "1.", m: New(map[string]interface{}{"a": 1, "b": 2}), callback: func(key string, val interface{}) bool { return false }},
		{name: "2.", m: New(map[string]interface{}{"a": 1, "b": 2}), callback: func(key string, val interface{}) bool { return true }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Iter(tt.callback)
		})
	}
}

func TestIgcMap_Keys(t *testing.T) {
	tests := []struct {
		name string
		m    *IgcMap
		want []string
	}{
		{name: "0.空map", m: New(map[string]interface{}{}), want: []string{}},
		{name: "1.一个key", m: New(map[string]interface{}{"a": 1}), want: []string{"a"}},
		{name: "2.多个key", m: New(map[string]interface{}{"a": 1, "b": 2}), want: []string{"a", "b"}},
		{name: "3.大小写key", m: New(map[string]interface{}{"a": 1, "A": "2"}), want: []string{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IgcMap.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIgcMap_Del(t *testing.T) {
	tests := []struct {
		name   string
		m      *IgcMap
		key    string
		keyLen int
	}{
		{name: "1.一个不存在的key-a", m: New(map[string]interface{}{}), key: "a", keyLen: 0},
		{name: "2.一个不存在的key-b", m: New(map[string]interface{}{"a": 1}), key: "b", keyLen: 1},
		{name: "3.存在的key", m: New(map[string]interface{}{"a": 1}), key: "a", keyLen: 0},
		{name: "4.大小写不一致", m: New(map[string]interface{}{"A": 1}), key: "a", keyLen: 0},
		{name: "5.多个key", m: New(map[string]interface{}{"a": 1, "b": 2}), key: "a", keyLen: 1},
		{name: "6.多个key-大小写", m: New(map[string]interface{}{"a": 1, "b": 2}), key: "A", keyLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Del(tt.key)
			if len(tt.m.data) != tt.keyLen {
				t.Errorf("IgcMap.Del() = %v, want %v", len(tt.m.data), tt.keyLen)
			}
		})
	}
}

func TestIgcMap_Orignal(t *testing.T) {
	tests := []struct {
		name string
		m    *IgcMap
		want map[string]interface{}
	}{
		{name: "1.空map", m: New(map[string]interface{}{}), want: map[string]interface{}{}},
		{name: "2.存在数据", m: New(map[string]interface{}{"a": 1}), want: map[string]interface{}{"a": 1}},
		{name: "3.存在多个数据", m: New(map[string]interface{}{"a": 1, "b": 2}), want: map[string]interface{}{"a": 1, "b": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Orignal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IgcMap.Orignal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIgcMap_Equal(t *testing.T) {

	tests := []struct {
		name string
		m    *IgcMap
		o    *IgcMap
		want bool
	}{
		{name: "1.", m: New(nil), o: nil, want: false},
		{name: "2.key一致", m: New(nil), o: New(map[string]interface{}{}), want: true},
		{name: "2.key一致,值不相等", m: New(map[string]interface{}{"a": 1}), o: New(map[string]interface{}{"a": 2}), want: false},
		{name: "3.key长度一致,key一致", m: New(map[string]interface{}{"a": 1}), o: New(map[string]interface{}{"a": 1}), want: true},
		{name: "4.key长度一致,key不一致", m: New(map[string]interface{}{"a": 1}), o: New(map[string]interface{}{"B": 1}), want: false},
		{name: "5.key长度一致,key一致,大小写不一致", m: New(map[string]interface{}{"a": 1}), o: New(map[string]interface{}{"A": 1}), want: true},
		{name: "2.key不一致", m: New(nil), o: New(map[string]interface{}{"a": 1}), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Equal(tt.o); got != tt.want {
				t.Errorf("IgcMap.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
