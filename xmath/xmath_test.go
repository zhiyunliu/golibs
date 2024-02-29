package xmath

import (
	"testing"
)

func TestMax(t *testing.T) {

	got := Max([]int{0, 1, 2}...)
	if got != 2 {
		t.Errorf("Max() = %v, want %v", got, 2)
	}

}
