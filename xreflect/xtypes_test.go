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
