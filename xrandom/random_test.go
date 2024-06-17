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

func TestRange(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
		want int
	}{
		{name: "1", min: 1, max: 10, want: 0},
		{name: "2", min: 1, max: 1, want: 1},
		{name: "2", min: 1, max: 0, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Range(tt.min, tt.max)

			if tt.want == 0 && !(got >= tt.min && got < tt.max) {
				t.Errorf("name=%s Range() = %v, ", tt.name, got)
			}

			if tt.want != 0 {
				if tt.want != got {
					t.Errorf("name=%s Range() = %v, want %v", tt.name, got, tt.want)
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
