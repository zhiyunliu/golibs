package xtypes

import (
	"testing"
)

func TestXMaps_Scan_1(t *testing.T) {
	type Args struct {
		A string
		B int
	}
	tests := []struct {
		name    string
		ms      XMaps
		args    []Args
		wantErr bool
	}{
		{name: "1.", ms: XMaps{XMap{"A": "a1", "B": 1}}, args: []Args{}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ms.Scan(&tt.args); (err != nil) != tt.wantErr {
				t.Errorf("XMaps.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.args) != 1 {
				t.Error("反射失败")
			}
		})
	}
}

func TestXMaps_Scan_2(t *testing.T) {
	type Args struct {
		A string
		B int
	}
	tests := []struct {
		name    string
		ms      XMaps
		args    []*Args
		wantErr bool
	}{
		{name: "1.", ms: XMaps{XMap{"A": "a2", "B": 2}}, args: []*Args{}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ms.Scan(&tt.args); (err != nil) != tt.wantErr {
				t.Errorf("XMaps.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.args) != 1 {
				t.Error("反射失败")
			}
		})
	}
}
