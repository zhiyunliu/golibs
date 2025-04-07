package httputil

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhiyunliu/golibs/xsse"
)

// MockResponseHandler 模拟响应处理器
func MockResponseHandler(resp *http.Response) (Body, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &normalBody{BodyBytes: body}, nil
}

// MockClient 模拟 HTTP 客户端
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		data           []byte
		opts           []Option
		respHandler    RespHandler
		sseHandler     SSEHandler
		expectedMethod string
		expectedURL    string
		expectedBody   string
		expectedError  string
	}{
		{
			name:           "ValidRequest",
			method:         "get",
			url:            "/aaa",
			data:           []byte(`{"key":"value"}`),
			opts:           []Option{WithHeader("Content-Type", "application/json")},
			expectedMethod: "GET",
			expectedURL:    "/aaa",
			expectedBody:   `{"key":"value"}`,
		},
		{
			name:           "RequestWithError",
			method:         "post",
			url:            "/aaa",
			data:           []byte(`{"key":"value"}`),
			opts:           []Option{WithHeader("Content-Type", "application/json")},
			expectedMethod: "POST",
			expectedURL:    "/aaa",
			expectedError:  "request error",
		},
		{
			name:           "RequestWithTLS",
			method:         "put",
			url:            "/aaa",
			data:           []byte(`{"key":"value"}`),
			opts:           []Option{WithTLS(&tls.Config{InsecureSkipVerify: true}), WithContentTypeJson()},
			respHandler:    MockResponseHandler,
			expectedMethod: "PUT",
			expectedURL:    "/aaa",
			expectedBody:   `{"key":"value"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 设置模拟服务器
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedMethod, r.Method)
				assert.Equal(t, test.expectedURL, r.URL.String())
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBody, string(body))
				if test.expectedError != "" {
					http.Error(w, test.expectedError, http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}))
			defer server.Close()

			// 设置模拟客户端
			mockClient := &MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if test.expectedError != "" {
						return nil, errors.New(test.expectedError)
					}
					return server.Client().Do(req)
				},
			}

			// 设置选项
			opts := append(test.opts, WithClient(mockClient))

			if test.respHandler != nil {
				opts = append(opts, WithRespHandler(test.respHandler))
			}

			// 调用被测试的方法
			body, err := Request(test.method, fmt.Sprintf("%s%s", server.URL, test.url), test.data, opts...)

			if test.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, body)
				body.GetHeader()
				body.GetResult()
				body.GetStatus()
			}
		})
	}
}

func TestSSERequest_normal(t *testing.T) {
	// 设置模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("wait for event")
		idx := 0
		for {
			idx++
			if idx > 3 {
				break
			}

			event := &xsse.Event{
				Id:    fmt.Sprint(idx),
				Event: "message",
				Data: map[string]any{
					"id": idx,
				},
			}
			err := xsse.HttpEncode(w, event)
			if err != nil {
				t.Log("decode event error:", err)
				return
			}
			w.(http.Flusher).Flush()
		}
		t.Log("finish  for event")

	}))
	defer server.Close()

	// 设置模拟客户端
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return server.Client().Do(req)
		},
	}
	url := fmt.Sprintf("%s/aaa", server.URL)
	data := []byte(`{"key":"value"}`)

	type Item struct {
		Id int `json:"id"`
	}

	unmarshaler := func(bytes []byte) (data any, err error) {
		var item Item
		err = json.Unmarshal(bytes, &item)
		data = item
		return
	}

	sseHandler := func(event *xsse.Event) error {
		t.Log("event:", event)
		return nil
	}

	// 设置选项
	opts := []Option{WithTLS(&tls.Config{InsecureSkipVerify: true}), WithClient(mockClient), WithSSEHandler(sseHandler, WithSSEUnmarshal(unmarshaler))}
	// 调用被测试的方法
	body, err := Request(http.MethodPost, url, data, opts...)
	assert.NoError(t, err)
	assert.NotNil(t, body)
	assert.Equal(t, int32(http.StatusOK), body.GetStatus())
}
