package xlist

import (
	"testing"
)

func TestList_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		l    *List
		want bool
	}{
		{name: "empty", l: NewList(), want: true},
		{name: "not empty", l: NewList().Append(1), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.IsEmpty(); got != tt.want {
				t.Errorf("List.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Length(t *testing.T) {
	tests := []struct {
		name string
		l    *List
		want int
	}{
		{name: "length=0", l: NewList(), want: 0},
		{name: "length=1", l: NewList().Append(1), want: 1},
		{name: "length=2", l: NewList().Append(1, 2), want: 2},
		{name: "length=3", l: NewList().Append(1, 2, 3), want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Length(); got != tt.want {
				t.Errorf("List.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Iter(t *testing.T) {
	type args struct {
		callback ListIterCallback
		rlength  int
	}
	tests := []struct {
		name string
		l    *List
		args args
	}{
		{
			name: "Iter=0", l: NewList(), args: args{
				rlength: 0,
				callback: func(idx int, node *Node) bool {
					t.Error("not in here")
					return true
				},
			},
		},
		{
			name: "Iter=1", l: NewList().Append(0), args: args{
				rlength: 1,
				callback: func(idx int, node *Node) bool {
					if !(idx == 0 && node.Value == idx) {
						t.Errorf("List.Iter = Iter=1")
						return false
					}
					return false
				},
			},
		},
		{
			name: "Iter=N", l: NewList().Append(0, 1, 2, 3, 4, 5, 6, 7, 8, 9), args: args{
				rlength: 10,
				callback: func(idx int, node *Node) bool {
					if idx != node.Value {
						t.Errorf("List.Iter-N = %d", idx)
						return false
					}
					return true
				},
			},
		},
		{
			name: "Iter=Remove", l: NewList().Append(0, 1, 2, 3, 4, 5, 6), args: args{
				rlength: 0,
				callback: func(idx int, node *Node) bool {
					node.Remove()
					return true
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Iter(tt.args.callback)
			if tt.args.rlength != tt.l.length {
				t.Errorf("Iter错误：%s", tt.name)
			}
		})
	}
}
