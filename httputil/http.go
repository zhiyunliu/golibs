package httputil

import (
	"strings"
	"sync"
)

const (
	_baseContentType       = "application"
	_contentType           = "Content-Type"
	_contentTypeJson       = "application/json;charset=utf-8"
	_contentTypeUrlencoded = "application/x-www-form-urlencoded;charset=utf-8"
	_contentTypeFormdata   = "multipart/form-data;charset=utf-8; boundary="
)

var _subTypeMap sync.Map

func init() {
	_subTypeMap.Store("text", "text/plain;charset=utf-8")
}

func ResetContentType(subtype, fullType string) {
	_subTypeMap.Store(subtype, fullType)
}

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	val, ok := _subTypeMap.Load(subtype)
	if ok {
		return val.(string)
	}
	return strings.Join([]string{_baseContentType, subtype}, "/")
}

// ContentSubtype returns the content-subtype for the given content-type.  The
// given content-type must be a valid content-type that starts with
// but no content-subtype will be returned.
// according rfc7231.
// contentType is assumed to be lowercase already.
func ContentSubtype(contentType string) string {

	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}
	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}
