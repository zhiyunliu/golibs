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
		{name: "empty", s: "", want: false},
		{name: "special_chars", s: "field-name", want: true},
		{name: "with_dot", s: "field.name", want: true},
		{name: "with_space", s: "field name", want: true},
		{name: "unicode_letter", s: "字段", want: true},
		{name: "digit_start", s: "1field", want: true},
		{name: "backslash", s: "field\\name", want: false},
		{name: "quote", s: "field\"name", want: false},
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

func TestParseTag_empty(t *testing.T) {
	name, opts := parseTag("")
	if name != "" {
		t.Errorf("parseTag(\"\") name = %q, want empty", name)
	}
	if opts != "" {
		t.Errorf("parseTag(\"\") opts = %q, want empty", opts)
	}
}

func TestParseTag_nameOnly(t *testing.T) {
	name, opts := parseTag("myfield")
	if name != "myfield" {
		t.Errorf("parseTag name = %q, want myfield", name)
	}
	if opts != "" {
		t.Errorf("parseTag opts = %q, want empty", opts)
	}
}

func TestTagOptions_Contains_empty(t *testing.T) {
	opts := TagOptions("")
	if opts.Contains("anything") {
		t.Error("empty options should not contain anything")
	}
}

func TestTagOptions_GetArgsInfo_empty(t *testing.T) {
	opts := TagOptions("")
	args, ok := opts.GetArgsInfo("anything")
	if ok {
		t.Error("empty options should not find anything")
	}
	if args != nil {
		t.Errorf("empty options args = %v, want nil", args)
	}
}

func TestTagOptions_GetArgsInfo_varchar(t *testing.T) {
	_, opts := parseTag("field,dbtype:varchar")
	args, ok := opts.GetArgsInfo("dbtype")
	if !ok {
		t.Error("GetArgsInfo(dbtype) should return true")
	}
	if len(args) != 1 || args[0] != "varchar" {
		t.Errorf("GetArgsInfo(dbtype) args = %v, want [varchar]", args)
	}
}

func TestTagOptions_GetArgsInfo_multipleArgs(t *testing.T) {
	_, opts := parseTag("field,dbtype:key1=val1&key2=val2")
	args, ok := opts.GetArgsInfo("dbtype")
	if !ok {
		t.Error("GetArgsInfo should return true")
	}
	expected := []string{"key1", "val1", "key2", "val2"}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("GetArgsInfo args = %v, want %v", args, expected)
	}
}
