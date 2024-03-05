package xtypes

import (
	"reflect"
	"testing"
)

func Test_mapscan(t *testing.T) {

	var val string = `{"a":1}`
	var bytes []byte = []byte(val)
	var eval string = `{"a":1`
	var ebytes []byte = []byte(val)
	var pval *string
	var pbytes *[]byte

	tests := []struct {
		name    string
		obj     any
		m       any
		wantErr bool
		expectM any
	}{
		{name: "1.", obj: nil, m: nil, wantErr: false},
		{name: "2.", obj: bytes, m: &XMap{}, wantErr: false, expectM: &XMap{"a": float64(1)}},
		{name: "3.", obj: val, m: &XMap{}, wantErr: false, expectM: &XMap{"a": float64(1)}},
		{name: "4.", obj: &bytes, m: &XMap{}, wantErr: false, expectM: &XMap{"a": float64(1)}},
		{name: "5.", obj: &val, m: &XMap{}, wantErr: false, expectM: &XMap{"a": float64(1)}},
		{name: "6.", obj: pval, m: &XMap{}, wantErr: false, expectM: &XMap{}},
		{name: "7.", obj: pbytes, m: &XMap{}, wantErr: false, expectM: &XMap{}},
		{name: "8.", obj: eval, m: &XMap{}, wantErr: true, expectM: &XMap{}},
		{name: "9.", obj: ebytes, m: &XMap{}, wantErr: false, expectM: &XMap{"a": float64(1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mapscan(tt.obj, tt.m); (err != nil) != tt.wantErr {
				t.Errorf("mapscan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.expectM, tt.m) {
				t.Errorf("mapscan() DeepEqual  expect=%v, m=%v", tt.expectM, tt.m)
			}
		})
	}
}
