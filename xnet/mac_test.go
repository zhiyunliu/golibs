package xnet

import (
	"net"
	"testing"
)

func TestGetLocalMac(t *testing.T) {
	// Mock interfaces for testing
	tests := []struct {
		name           string
		mockInterfaces func() ([]net.Interface, error)
		expectedMac    string
	}{
		{
			name: "Happy Path",
			mockInterfaces: func() ([]net.Interface, error) {
				return []net.Interface{
					{Name: "eth0", HardwareAddr: []byte{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E}},
				}, nil
			},
			expectedMac: "00:1A:2B:3C:4D:5E",
		},
		{
			name: "No Interfaces",
			mockInterfaces: func() ([]net.Interface, error) {
				return []net.Interface{}, nil
			},
			expectedMac: StaticLocalMac,
		},
		{
			name: "Empty Hardware Address",
			mockInterfaces: func() ([]net.Interface, error) {
				return []net.Interface{
					{Name: "eth0", HardwareAddr: []byte{}},
				}, nil
			},
			expectedMac: StaticLocalMac,
		},
		{
			name: "Error Fetching Interfaces",
			mockInterfaces: func() ([]net.Interface, error) {
				return nil, net.UnknownNetworkError("mock error")
			},
			expectedMac: StaticLocalMac,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	GetInterfaces = tt.mockInterfaces
			result := GetLocalMac()
			if result != tt.expectedMac {
				t.Errorf("GetLocalMac() = %v, want %v", result, tt.expectedMac)
			}
		})
	}
}
