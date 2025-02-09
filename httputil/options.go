package httputil

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
)

// RespHandler is a function that takes a io.Reader and returns a io.Reader
type RespHandler func(resp *http.Response) (Body, error)

type options struct {
	header      http.Header
	client      Client
	tls         *tls.Config
	respHandler RespHandler
}

func defaultOptions() *options {
	return &options{
		client: http.DefaultClient,
		respHandler: func(resp *http.Response) (Body, error) {

			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			header := make(map[string]string)
			for k, v := range resp.Header {
				header[k] = strings.Join(v, ",")
			}

			body := &normalBody{
				Status: int32(resp.StatusCode),
				Body:   respBytes,
				Header: header,
			}

			return body, nil
		},
	}
}

// Option is a function that takes a *options and modifies it
type Option func(o *options)

// WithHeader sets the header of the request
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

// WithContentTypeJson sets the content type of the request to json
func WithContentTypeJson() Option {
	return WithContentType(_contentTypeJson)
}

// WithContentTypeUrlencoded sets the content type of the request to urlencoded
func WithContentTypeUrlencoded() Option {
	return WithContentType(_contentTypeUrlencoded)
}

// WithContentTypeFormData sets the content type of the request to form data
func WithContentTypeFormData(boundary string) Option {
	if strings.ContainsAny(boundary, `()<>@,;:\"/[]?= `) {
		boundary = `"` + boundary + `"`
	}
	return WithContentType(_contentTypeFormdata + boundary)
}

// WithContentType sets the content type of the request
func WithContentType(contentType string) Option {
	return func(o *options) {
		if o.header == nil {
			o.header = make(http.Header)
		}
		o.header[_contentType] = []string{contentType}
	}
}

// WithClient sets the http client
func WithClient(client Client) Option {
	return func(o *options) {
		o.client = client
	}
}

// WithTLS sets the tls connection state
func WithTLS(tlscfg *tls.Config) Option {
	return func(o *options) {
		o.tls = tlscfg
	}
}

// BodyHandler is a function that takes a io.Reader and returns a io.Reader
func WithRespHandler(handler RespHandler) Option {
	return func(o *options) {
		o.respHandler = handler
	}
}
