// Copyright 2011 The Go Authors. All rights reserved.

// Use of this source code is governed by a BSD-style

// license that can be found in the LICENSE file.

package xreflect

import (
	"reflect"
	"strings"
	"testing"
)

func TestTagParsing(t *testing.T) {
	name, opts := parseTag("field,foobar,foo,dbtype:vtp=utp_sup_stock_item")
	if name != "field" {
		t.Fatalf("name = %q, want field", name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"foobar", true},
		{"foo", true},
		{"bar", false},
		{"dbtype", true},
	} {
		if opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%q) = %v", tt.opt, !tt.want)
		}
	}
}

func Test_isValidTag(t *testing.T) {

	tests := []struct {
		name string
		s    string
		want bool
	}{
		{name: "1", s: "field", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidTag(tt.s); got != tt.want {
				t.Errorf("isValidTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagOptions_GetArgInfo(t *testing.T) {

	name, opts := parseTag("field,foobar,foo,dbtype:vtp=utp_sup_stock_item")
	if name != "field" {
		t.Fatalf("name = %q, want field", name)
	}
	for _, tt := range []struct {
		opt     string
		argInfo []string
		want    bool
	}{
		{opt: "foobar", argInfo: nil, want: true},
		{opt: "foo", argInfo: nil, want: true},
		{opt: "bar", argInfo: nil, want: false},
		{opt: "dbtype", argInfo: []string{"vtp", "utp_sup_stock_item"}, want: true},
	} {
		argInfo, ok := opts.GetArgsInfo(tt.opt)
		if ok != tt.want {
			t.Errorf("GetArgInfo(%q) = %v", tt.opt, !tt.want)
		}
		if !reflect.DeepEqual(argInfo, tt.argInfo) {
			t.Errorf("GetArgInfo(%q) ,Got:%v,want:%s", tt.opt, argInfo, tt.argInfo)
		}
	}

}

func TestXxx(t *testing.T) {
	tmparg := "aaaa=bbb&ccc=ddd"

	argInfo := strings.FieldsFunc(tmparg, func(r rune) bool {
		if r == '=' || r == '&' {
			return true
		}
		return false
	})

	t.Log(argInfo)

}
