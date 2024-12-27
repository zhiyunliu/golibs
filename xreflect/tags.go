// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xreflect

import (
	"strings"
	"unicode"
)

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type TagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, TagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, TagOptions(opt)
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o TagOptions) Contains(optionName string) bool {
	//dbtype:vtp=utp_sup_stock_item
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
		if !strings.Contains(name, ":") {
			continue
		}

		name, _, _ = strings.Cut(name, ":")
		if name == optionName {
			return true
		}
	}
	return false
}

func (o TagOptions) GetArgsInfo(opt string) (args []string, ok bool) {
	//dbtype:vtp=utp_sup_stock_item  ==>
	//dbtype:varchar
	if len(o) == 0 {
		return args, false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == opt {
			return args, true
		}
		if !strings.Contains(name, ":") {
			continue
		}

		name, tmparg, _ := strings.Cut(name, ":")
		if name == opt {
			args = strings.FieldsFunc(tmparg, func(r rune) bool {
				if r == '=' || r == '&' {
					return true
				}
				return false
			})
			return args, true
		}
	}
	return args, false
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}
