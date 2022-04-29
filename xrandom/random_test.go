package xrandom

import (
	"testing"
)

func BenchmarkRange(t *testing.B) {

	tests := []struct {
		name string
		min  int
		max  int
	}{
		{name: "1.", min: 1, max: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			for i := 0; i < t.N; i++ {
				got := Range(tt.min, tt.max)
				if !(tt.min <= got && got < tt.max) {
					t.Errorf("Range() = %v, min=%d,max=%d", got, tt.min, tt.max)
				}
			}
		})
	}
}

func TestStr(t *testing.T) {

	tests := []struct {
		name    string
		n       int
		wantLen int
	}{
		{name: "1", n: 1, wantLen: 1},
		{name: "2", n: 10, wantLen: 10},
		{name: "3", n: 100, wantLen: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Str(tt.n); len(got) != tt.wantLen {
				t.Errorf("Str() = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func BenchmarkStr(t *testing.B) {

	tests := []struct {
		name string
		xlen int
	}{
		{name: "1.", xlen: 1},
		{name: "2.", xlen: 100},
		{name: "3.", xlen: 10000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			for i := 0; i < t.N; i++ {
				got := Str(tt.xlen)
				if tt.xlen != len(got) {
					t.Errorf("name() = %v, xlen=%d", tt.name, tt.xlen)
				}
			}
		})
	}
}
