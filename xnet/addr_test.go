package xnet

import (
	"strconv"
	"strings"
	"testing"

	"github.com/zhiyunliu/golibs/xlog"
)

type testLogger struct {
}

func (t testLogger) Name() string {
	return "test"
}
func (t testLogger) SessionID() string {
	return "testid"
}
func (t testLogger) Log(level xlog.Level, args ...interface{}) {

}
func (t testLogger) Logf(level xlog.Level, format string, args ...interface{}) {

}
func (t testLogger) Close() {

}

func TestGetAvaliableAddr(t *testing.T) {
	logger := &testLogger{}
	localIp := "127.0.0.1"
	var min int64 = 1000
	var max int64 = 2000
	tests := []struct {
		name        string
		addr        string
		min         int64
		max         int64
		wantNewAddr string
		wantErr     bool
	}{
		{name: "1.", addr: ":1001", wantNewAddr: ":1001", wantErr: false},
		{name: "2.", addr: "[1000,2000)", wantNewAddr: "", wantErr: false},
		{name: "3.", addr: "[1000,2000):rand", wantNewAddr: "", wantErr: false},
		{name: "4.", addr: "[1000,2000):seq:5", wantNewAddr: "", wantErr: false},
		{name: "5.", addr: "[1000,2000):xseq:5", wantNewAddr: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewAddr, err := GetAvaliableAddr(logger, localIp, tt.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvaliableAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				return
			}

			parties := strings.Split(gotNewAddr, ":")
			nport, err := strconv.ParseInt(parties[1], 10, 64)
			if err != nil {
				t.Error(err)
			}
			if !(min <= nport && nport < max) {
				t.Error(err)
			}
			if tt.wantNewAddr != "" && gotNewAddr != tt.wantNewAddr {
				t.Errorf("GetAvaliableAddr() = %v, want %v", gotNewAddr, tt.wantNewAddr)
			}
		})
	}
}
