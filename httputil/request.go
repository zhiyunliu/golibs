package httputil

import (
	"bytes"
	"net/http"
	"strings"
)

type Body interface {
	GetStatus() int32
	GetHeader() map[string]string
	GetResult() []byte
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type normalBody struct {
	Status int32
	Header map[string]string
	Body   []byte
}

func (b *normalBody) GetStatus() int32 {
	return b.Status
}
func (b *normalBody) GetHeader() map[string]string {
	return b.Header
}
func (b *normalBody) GetResult() []byte {
	return b.Body
}

// Request sends a request to the given URL with the given method and data.
func Request(method string, url string, data []byte, opts ...Option) (body Body, err error) {
	opt := defaultOptions()

	for i := range opts {
		opts[i](opt)
	}
	client := opt.client
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
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

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = opt.respHandler(resp)

	return body, err
}
