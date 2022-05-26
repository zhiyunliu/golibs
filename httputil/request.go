package httputil

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

type Body interface {
	GetStatus() int32
	GetHeader() map[string]string
	GetResult() []byte
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

func Request(method string, url string, data []byte, opts ...Option) (body Body, err error) {
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
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	header := make(map[string]string)
	for k, v := range resp.Header {
		header[k] = strings.Join(v, ",")
	}

	body = &normalBody{
		Status: int32(resp.StatusCode),
		Body:   respBytes,
		Header: header,
	}

	return body, err
}
