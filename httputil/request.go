package httputil

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/zhiyunliu/golibs/xsse"
)

var (
	_ Body           = (*normalBody)(nil)
	_ ResponseEntity = (*normalBody)(nil)
)

type ResponseEntity interface {
	StatusCode() int
	Header() map[string]string
	Body() (bytes []byte, err error)
}

type Body interface {
	GetStatus() int32
	GetHeader() map[string]string
	GetResult() []byte
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type normalBody struct {
	Status    int32
	Headers   map[string]string
	BodyBytes []byte
}

func (b *normalBody) GetStatus() int32 {
	return b.Status
}
func (b *normalBody) GetHeader() map[string]string {
	return b.Headers
}
func (b *normalBody) GetResult() []byte {
	return b.BodyBytes
}

func (b *normalBody) StatusCode() int {
	return int(b.Status)
}
func (b *normalBody) Header() map[string]string {
	return b.Headers
}
func (b *normalBody) Body() (bytes []byte, err error) {
	return b.BodyBytes, nil
}

func NewEmptyBody() Body {
	return &normalBody{
		Status:    http.StatusOK,
		Headers:   make(map[string]string),
		BodyBytes: make([]byte, 0),
	}
}

// RequestWithContext sends a request to the given URL with the given method and data.
func RequestWithContext(ctx context.Context, method string, url string, data []byte, opts ...Option) (body Body, err error) {
	opt := defaultOptions()

	for i := range opts {
		opts[i](opt)
	}
	client := opt.client
	method = strings.ToUpper(method)
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header = opt.header
	if opt.tls != nil {
		transport := &http.Transport{
			TLSClientConfig: opt.tls,
		}
		client = &http.Client{
			Transport: transport,
		}
	}
	opt.ReqChangeCalls.Apply(req)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if opt.sseHandler != nil {
		return handleSSE(opt, resp.Body)
	}
	return opt.respHandler(resp)
}

// Request sends a request to the given URL with the given method and data.
func Request(method string, url string, data []byte, opts ...Option) (body Body, err error) {
	return RequestWithContext(context.Background(), method, url, data, opts...)
}

func handleSSE(opt *options, respBody io.Reader) (body Body, err error) {
	body = NewEmptyBody()
	handler := opt.sseHandler
	events, err := xsse.Decode(respBody, opt.sseOpts...)
	if err != nil {
		return
	}
	for event := range events {
		if err = handler(event); err != nil {
			return
		}
	}
	return
}
