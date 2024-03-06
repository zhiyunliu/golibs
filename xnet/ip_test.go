package xnet

import (
	"net"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-playground/assert/v2"
)

func TestGetLocalIP(t *testing.T) {

	patches := gomonkey.ApplyFunc(getInterfaceAddrs, func() ([]net.Addr, error) {
		var result []net.Addr
		result = append(result, &net.IPNet{IP: net.ParseIP("192.168.1.100")})
		return result, nil

	})
	// 每个单测过后需要reset
	defer patches.Reset()

	localIp := GetLocalIP("192.168")

	assert.Equal(t, localIp, "192.168.1.100")
}

func TestExtractHostPort(t *testing.T) {
	// Testing with a valid address string
	validAddr := "example.com:8080"
	host, port, err := ExtractHostPort(validAddr)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if host != "example.com" || port != 8080 {
		t.Errorf("Expected host:example.com, port:8080, but got host:%s, port:%d", host, port)
	}

	// Testing with an invalid address string (missing port)
	invalidAddr1 := "example.com:"
	_, _, err = ExtractHostPort(invalidAddr1)
	if err == nil {
		t.Error("Expected an error for invalid address, but no error was returned")
	}

	// Testing with an invalid address string (wrong format)
	invalidAddr2 := "localhost"
	_, _, err = ExtractHostPort(invalidAddr2)
	if err == nil {
		t.Error("Expected an error for invalid address format, but no error was returned")
	}
}

func TestIsValidIP(t *testing.T) {
	// Testing with a valid global unicast IP address
	validIP := "192.0.2.1"
	result := IsValidIP(validIP)
	if !result {
		t.Errorf("Expected true for valid global unicast IP, but got false")
	}

	// Testing with an invalid IP address (interface-local multicast)
	invalidIP := "ff01::1"
	result = IsValidIP(invalidIP)
	if result {
		t.Errorf("Expected false for interface-local multicast IP, but got true")
	}
}

func Test_getInterfaceAddrs(t *testing.T) {

	patches := gomonkey.ApplyFunc(net.InterfaceAddrs, func() ([]net.Addr, error) {
		var result []net.Addr
		result = append(result, &net.IPNet{IP: net.ParseIP("192.168.1.100")})
		return result, nil

	})
	// 每个单测过后需要reset
	defer patches.Reset()

	gotAddrs, err := getInterfaceAddrs()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(gotAddrs) != 1 {

		t.Errorf("Expected 1 address, but got: %d", len(gotAddrs))
	}

}
