package httputil

import (
	"crypto/tls"
	"net/http"
	"strings"
)

type options struct {
	header http.Header
	client *http.Client
	tls    *tls.ConnectionState
}

type Option func(o *options)

func WithHeader(name string, val ...string) Option {
	return func(o *options) {
		if o.header == nil {
			o.header = make(http.Header)
		}
		if len(name) == 0 || len(val) == 0 {
			return
		}
		o.header[name] = val
	}
}

func WithContentTypeJson() Option {
	return WithContentType(_contentTypeJson)
}

func WithContentTypeUrlencoded() Option {
	return WithContentType(_contentTypeUrlencoded)
}

func WithContentTypeFormData(boundary string) Option {
	if strings.ContainsAny(boundary, `()<>@,;:\"/[]?= `) {
		boundary = `"` + boundary + `"`
	}
	return WithContentType(_contentTypeFormdata + boundary)
}

func WithContentType(contentType string) Option {
	return func(o *options) {
		if o.header == nil {
			o.header = make(http.Header)
		}
		o.header[_contentType] = []string{contentType}
	}
}

func WithClient(client *http.Client) Option {
	return func(o *options) {
		o.client = client
	}
}

func WithTLS(tlsCert *tls.ConnectionState) Option {
	return func(o *options) {
		o.tls = tlsCert
	}
}
