package xnet

import "testing"

func TestParse(t *testing.T) {
	// Happy path test case
	t.Run("Valid address", func(t *testing.T) {
		wantProto := "https"
		wantName := "example.com"
		addr := "  https://example.com  "
		gotProto, gotName, err := Parse(addr)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if gotProto != wantProto || gotName != wantName {
			t.Errorf("Parse(%q) = (%q, %q, nil), want (%q, %q, nil)", addr, gotProto, gotName, wantProto, wantName)
		}
	})

}
