package session

import (
	"testing"
)

func TestCreateSession(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "1", want: "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Create(); got != tt.want {
				t.Errorf("CreateSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
