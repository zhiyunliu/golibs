package xreflect

import (
	"math"
	"testing"
)

type S1 struct {
	S string
}

type S2 struct {
	S string
}

func (s S2) String() string {
	return s.S
}

func TestGetString(t *testing.T) {

	ptrBytes := []byte("123")
	ptrStr := "123"

	tests := []struct {
		name string
		v    interface{}
		want string
	}{
		{name: "1.", v: 1, want: "1"},
		{name: "2.", v: 1.0, want: "1"},
		{name: "3.", v: 1.1, want: "1.1"},
		{name: "4.", v: "a", want: "a"},
		{name: "5.", v: S1{S: "s"}, want: "{S:s}"},
		{name: "6.", v: S2{S: "s"}, want: "s"},
		{name: "7.", v: nil, want: ""},
		{name: "8.", v: true, want: "true"},
		{name: "9.", v: false, want: "false"},
		{name: "10.", v: 10000000000, want: "10000000000"},
		{name: "11.", v: uint64(10000000000), want: "10000000000"},
		{name: "12.", v: float64(10000000000.10), want: "10000000000.1"},
		{name: "13.", v: []byte("aaa"), want: "aaa"},
		{name: "14.", v: &ptrBytes, want: "123"},
		{name: "15.", v: &ptrStr, want: "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetString(tt.v); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestGetBool(t *testing.T) {

	tests := []struct {
		name string
		tmp  interface{}

		want bool
	}{
		{name: "1.", tmp: nil, want: false},
		{name: "2.", tmp: "1", want: true},
		{name: "3.", tmp: "0", want: false},
		{name: "4.", tmp: "aaa", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBool(tt.tmp); got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {

	var val int = 1
	var val8 int8 = 1
	var val16 int16 = 1
	var val32 int32 = 1
	var val64 int64 = math.MaxInt64
	var val64s int64 = 2
	tests := []struct {
		name    string
		tmp     interface{}
		want    int
		wantErr bool
	}{
		{name: "1.", tmp: nil, want: 0, wantErr: false},
		{name: "2.", tmp: val, want: 1, wantErr: false},
		{name: "3.", tmp: &val, want: 1, wantErr: false},

		{name: "4.", tmp: val8, want: 1, wantErr: false},
		{name: "5.", tmp: &val8, want: 1, wantErr: false},

		{name: "6.", tmp: val16, want: 1, wantErr: false},
		{name: "7.", tmp: &val16, want: 1, wantErr: false},

		{name: "8.", tmp: val32, want: 1, wantErr: false},
		{name: "9.", tmp: &val32, want: 1, wantErr: false},

		{name: "10.", tmp: val64, want: math.MaxInt, wantErr: false},
		{name: "11.", tmp: &val64, want: math.MaxInt, wantErr: false},

		{name: "12.", tmp: val64s, want: 2, wantErr: false},
		{name: "13.", tmp: &val64s, want: 2, wantErr: false},

		{name: "14.", tmp: "1", want: 1, wantErr: false},
		{name: "15.", tmp: S2{S: "1"}, want: 1, wantErr: false},
		{name: "16.", tmp: S1{S: "1"}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInt(tt.tmp)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt64(t *testing.T) {
	var (
		val   int64 = 1
		val8  int8  = 1
		val16 int16 = 1
		val32 int32 = 1
	)

	tests := []struct {
		name    string
		tmp     interface{}
		want    int64
		wantErr bool
	}{
		{name: "1.", tmp: nil, want: 0, wantErr: false},
		{name: "2.", tmp: val, want: 1, wantErr: false},
		{name: "3.", tmp: &val, want: 1, wantErr: false},
		{name: "4.", tmp: val8, want: 1, wantErr: false},
		{name: "5.", tmp: &val8, want: 1, wantErr: false},
		{name: "6.", tmp: val16, want: 1, wantErr: false},
		{name: "7.", tmp: &val16, want: 1, wantErr: false},
		{name: "8.", tmp: val32, want: 1, wantErr: false},
		{name: "9.", tmp: &val32, want: 1, wantErr: false},
		{name: "10.", tmp: "1", want: 1, wantErr: false},
		{name: "11.", tmp: S2{S: "1"}, want: 1, wantErr: false},
		{name: "12.", tmp: S1{S: "1"}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInt64(tt.tmp)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUint64(t *testing.T) {
	var (
		uval   uint64 = 1
		uval8  uint8  = 1
		uval16 uint16 = 1
		uval32 uint32 = 1
		ival   int    = 1
		ival8  int8   = 1
		ival16 int16  = 1
		ival32 int32  = 1
		ival64 int64  = 1
	)

	tests := []struct {
		name    string
		tmp     interface{}
		want    uint64
		wantErr bool
	}{
		{name: "1.", tmp: nil, want: 0, wantErr: false},
		{name: "2.", tmp: uval, want: 1, wantErr: false},
		{name: "3.", tmp: &uval, want: 1, wantErr: false},
		{name: "4.", tmp: uval8, want: 1, wantErr: false},
		{name: "5.", tmp: &uval8, want: 1, wantErr: false},
		{name: "6.", tmp: uval16, want: 1, wantErr: false},
		{name: "7.", tmp: &uval16, want: 1, wantErr: false},
		{name: "8.", tmp: uval32, want: 1, wantErr: false},
		{name: "9.", tmp: &uval32, want: 1, wantErr: false},
		{name: "10.", tmp: ival, want: 1, wantErr: false},
		{name: "11.", tmp: &ival, want: 1, wantErr: false},
		{name: "12.", tmp: ival8, want: 1, wantErr: false},
		{name: "13.", tmp: &ival8, want: 1, wantErr: false},
		{name: "14.", tmp: ival16, want: 1, wantErr: false},
		{name: "15.", tmp: &ival16, want: 1, wantErr: false},
		{name: "16.", tmp: ival32, want: 1, wantErr: false},
		{name: "17.", tmp: &ival32, want: 1, wantErr: false},
		{name: "18.", tmp: ival64, want: 1, wantErr: false},
		{name: "19.", tmp: &ival64, want: 1, wantErr: false},
		{name: "20.", tmp: "1", want: 1, wantErr: false},
		{name: "21.", tmp: S2{S: "1"}, want: 1, wantErr: false},
		{name: "22.", tmp: S1{S: "1"}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUint64(tt.tmp)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFloat64(t *testing.T) {
	var (
		f32 float32 = 1.1
		f64 float64 = 1.1
	)

	// 辅助函数：判断两个float64是否近似相等
	almostEqual := func(a, b float64) bool {
		return math.Abs(a-b) < 1e-5
	}

	tests := []struct {
		name    string
		tmp     interface{}
		want    float64
		wantErr bool
	}{
		{name: "1.", tmp: nil, want: 0, wantErr: false},
		{name: "2.", tmp: f32, want: 1.1, wantErr: false},
		{name: "3.", tmp: &f32, want: 1.1, wantErr: false},
		{name: "4.", tmp: f64, want: 1.1, wantErr: false},
		{name: "5.", tmp: &f64, want: 1.1, wantErr: false},
		{name: "6.", tmp: "1.1", want: 1.1, wantErr: false},
		{name: "7.", tmp: S2{S: "1.1"}, want: 1.1, wantErr: false},
		{name: "8.", tmp: S1{S: "1.1"}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFloat64(tt.tmp)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 对于浮点数，使用近似相等判断
			if !almostEqual(got, tt.want) {
				t.Errorf("GetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strToint(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    int
		wantErr bool
	}{
		{name: "normal", str: "42", want: 42, wantErr: false},
		{name: "zero", str: "0", want: 0, wantErr: false},
		{name: "negative", str: "-1", want: -1, wantErr: false},
		{name: "invalid", str: "abc", want: 0, wantErr: true},
		{name: "empty", str: "", want: 0, wantErr: true},
		{name: "float_string", str: "1.5", want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToint(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("strToint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("strToint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newNotSupportErr(t *testing.T) {
	err := newNotSupportErr(42)
	if err == nil {
		t.Error("newNotSupportErr should return an error")
	}
	if err.Error() != "不支持的数据类型:int" {
		t.Errorf("newNotSupportErr error = %v", err)
	}
}

func TestGetInt_boundaryOverflow(t *testing.T) {
	// int64 max value that overflows int on 32-bit (but on 64-bit it's fine)
	// This test ensures the overflow check works
	var val64 int64 = math.MaxInt64
	got, _ := GetInt(val64)
	// On 64-bit systems, MaxInt64 == MaxInt, so no overflow
	if got != math.MaxInt {
		t.Errorf("GetInt(MaxInt64) = %v, want %v", got, math.MaxInt)
	}
}

func TestGetInt_pointerOverflow(t *testing.T) {
	var val64 int64 = math.MaxInt64
	got, _ := GetInt(&val64)
	if got != math.MaxInt {
		t.Errorf("GetInt(&MaxInt64) = %v, want %v", got, math.MaxInt)
	}
}

func TestGetString_numTypes(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want string
	}{
		{name: "int8", v: int8(8), want: "8"},
		{name: "int16", v: int16(16), want: "16"},
		{name: "int32", v: int32(32), want: "32"},
		{name: "int64", v: int64(64), want: "64"},
		{name: "uint8", v: uint8(8), want: "8"},
		{name: "uint16", v: uint16(16), want: "16"},
		{name: "uint32", v: uint32(32), want: "32"},
		{name: "float32", v: float32(1.5), want: "1.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetString(tt.v); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt64_withInt(t *testing.T) {
	var v int = 42
	got, err := GetInt64(v)
	if err != nil {
		t.Errorf("GetInt64(int) error = %v", err)
	}
	if got != 42 {
		t.Errorf("GetInt64(int) = %v, want 42", got)
	}
}

func TestGetUint64_withUint(t *testing.T) {
	var v uint = 42
	got, err := GetUint64(v)
	if err != nil {
		t.Errorf("GetUint64(uint) error = %v", err)
	}
	if got != 42 {
		t.Errorf("GetUint64(uint) = %v, want 42", got)
	}
}
