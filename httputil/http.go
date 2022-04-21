package httputil

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const (
	_baseContentType = "application"
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

func Request(method string, url string, data []byte, opts ...Option) (respBytes []byte, err error) {
	opt := &options{
		client: http.DefaultClient,
	}

	for i := range opts {
		opts[i](opt)
	}
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header = opt.header
	if opt.tls != nil {
		req.TLS = opt.tls
	}

	resp, err := opt.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBytes, err = ioutil.ReadAll(resp.Body)
	return
}
