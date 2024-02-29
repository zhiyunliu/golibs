package xtypes

import (
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
