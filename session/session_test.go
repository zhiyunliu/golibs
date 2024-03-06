package session

import (
	"testing"
)

func TestCreateSession(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
	}{
		{name: "1", wantLen: 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Create(); len(got) != tt.wantLen {
				t.Errorf("CreateSession() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}
