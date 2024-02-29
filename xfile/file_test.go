package xfile

import (
	"path/filepath"
	"strings"
	"testing"
)

func Test_getAbsFilePath(t *testing.T) {
	basePath, err := filepath.Abs(".")
	if err != nil {
		t.Log(err)
	}
	tests := []struct {
		name        string
		path        string
		wantAbsPath string
		wantErr     bool
	}{
		{name: "1.", path: "../log/a.log", wantAbsPath: filepath.Dir(basePath) + "/log/a.log", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAbsPath, err := getAbsFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAbsFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotAbsPath = strings.ReplaceAll(gotAbsPath, string(filepath.Separator), "/")
			tt.wantAbsPath = strings.ReplaceAll(tt.wantAbsPath, string(filepath.Separator), "/")
			if gotAbsPath != tt.wantAbsPath {
				t.Errorf("getAbsFilePath() = %v, want %v", gotAbsPath, tt.wantAbsPath)
			}
		})
	}
}
