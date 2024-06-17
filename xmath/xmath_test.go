package xmath

import (
	"reflect"
	"testing"
)

func TestMax(t *testing.T) {

	got := Max([]int{0, 1, 2}...)
	if got != 2 {
		t.Errorf("Max() = %v, want %v", got, 2)
	}

	got1 := Max([]int{}...)
	if got1 != 0 {
		t.Errorf("Max() = %v, want %v", got1, 0)
	}

}
func TestMin(t *testing.T) {

	got := Min([]int{1, 0, 2}...)
	if got != 0 {
		t.Errorf("Min() = %v, want %v", got, 0)
	}
	got1 := Min([]int{}...)
	if got1 != 0 {
		t.Errorf("Min() = %v, want %v", got1, 0)
	}

}

func TestSum(t *testing.T) {

	got := Sum([]int{0, 1, 2}...)
	if got != 3 {
		t.Errorf("Sum() = %v, want %v", got, 3)
	}

}

func TestAbs(t *testing.T) {
	val1 := -1
	if got := Abs(val1); !reflect.DeepEqual(got, 1) {
		t.Errorf("Abs() = %v, want %v", got, 1)
	}

	val2 := 1
	if got := Abs(val2); !reflect.DeepEqual(got, 1) {
		t.Errorf("Abs() = %v, want %v", got, 1)
	}
}
