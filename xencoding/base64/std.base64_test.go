package base64

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		wantS   []byte
		wantErr bool
	}{
		{name: "test1", args: args{src: "aGVsbG8gd29ybGQ="}, wantS: []byte("hello world"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := Decode(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("Decode() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
