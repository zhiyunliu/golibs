package httputil

import (
	"strings"
	"testing"
)

func TestContentType(t *testing.T) {
	// Happy path
	_subTypeMap.Store("json", _contentTypeJson)
	_subTypeMap.Store("urlencoded", _contentTypeUrlencoded)
	_subTypeMap.Store("formdata", _contentTypeFormdata)

	if got := ContentType("json"); got != _contentTypeJson {
		t.Errorf("ContentType('json') = %v; want %v", got, _contentTypeJson)
	}

	if got := ContentType("urlencoded"); got != _contentTypeUrlencoded {
		t.Errorf("ContentType('urlencoded') = %v; want %v", got, _contentTypeUrlencoded)
	}

	if got := ContentType("formdata"); strings.HasPrefix(got, _contentTypeFormdata) == false {
		t.Errorf("ContentType('formdata') = %v; want %v", got, _contentTypeFormdata)
	}

	// Edge cases
	if got := ContentType("text"); got != "text/plain;charset=utf-8" {
		t.Errorf("ContentType('text') = %v; want %v", got, "text/plain;charset=utf-8")
	}

	if got := ContentType("invalid-subtype"); got != "application/invalid-subtype" {
		t.Errorf("ContentType('invalid-subtype') = %v; want %v", got, "application/invalid-subtype")
	}

	// Case sensitivity
	_subTypeMap.Store("JSON", "application/JSON;charset=utf-8")
	if got := ContentType("JSON"); got != "application/JSON;charset=utf-8" {
		t.Errorf("ContentType('JSON') = %v; want %v", got, "application/JSON;charset=utf-8")
	}

	if got := ContentType("json"); got != _contentTypeJson {
		t.Errorf("ContentType('json') = %v; want %v", got, _contentTypeJson)
	}
}

func TestContentSubtype(t *testing.T) {
	// Happy path
	t.Run("Valid JSON Content-Type", func(t *testing.T) {
		result := ContentSubtype("application/json;charset=utf-8")
		if result != "json" {
			t.Errorf("Expected 'json', got '%s'", result)
		}
	})
	t.Run("Valid URL Encoded Content-Type", func(t *testing.T) {
		result := ContentSubtype("application/x-www-form-urlencoded;charset=utf-8")
		if result != "x-www-form-urlencoded" {
			t.Errorf("Expected 'x-www-form-urlencoded', got '%s'", result)
		}
	})
	t.Run("Valid Custom Subtype", func(t *testing.T) {
		ResetContentType("custom", "application/custom;charset=utf-8")
		result := ContentSubtype("application/custom;charset=utf-8")
		if result != "custom" {
			t.Errorf("Expected 'custom', got '%s'", result)
		}
	})

	// Edge cases
	t.Run("Empty String", func(t *testing.T) {
		result := ContentSubtype("")
		if result != "" {
			t.Errorf("Expected '', got '%s'", result)
		}
	})
	t.Run("No Slash", func(t *testing.T) {
		result := ContentSubtype("applicationjson;charset=utf-8")
		if result != "" {
			t.Errorf("Expected '', got '%s'", result)
		}
	})
	t.Run("Slash Only", func(t *testing.T) {
		result := ContentSubtype("/")
		if result != "" {
			t.Errorf("Expected '', got '%s'", result)
		}
	})
	t.Run("No Semicolon", func(t *testing.T) {
		result := ContentSubtype("application/xml")
		if result != "xml" {
			t.Errorf("Expected 'xml', got '%s'", result)
		}
	})
	t.Run("Invalid Prefix", func(t *testing.T) {
		result := ContentSubtype("text/json;charset=utf-8")
		if result != "json" {
			t.Errorf("Expected 'json', got '%s'", result)
		}
	})

	t.Run("Semicolon Before Slash", func(t *testing.T) {
		result := ContentSubtype("application;foo/bar;charset=utf-8")
		if result != "" {
			t.Errorf("Expected '', got '%s'", result)
		}
	})
}
