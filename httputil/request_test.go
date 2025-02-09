package httputil

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockResponseHandler 模拟响应处理器
func MockResponseHandler(resp *http.Response) (Body, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &normalBody{Body: body}, nil
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
		expectedMethod string
		expectedURL    string
		expectedBody   string
		expectedError  string
	}{
		{
			name:           "ValidRequest",
			method:         "get",
			url:            "http://example.com/aaa",
			data:           []byte(`{"key":"value"}`),
			opts:           []Option{WithHeader("Content-Type", "application/json")},
			expectedMethod: "GET",
			expectedURL:    "/aaa",
			expectedBody:   `{"key":"value"}`,
		},
		{
			name:           "RequestWithError",
			method:         "post",
			url:            "http://example.com/aaa",
			data:           []byte(`{"key":"value"}`),
			opts:           []Option{WithHeader("Content-Type", "application/json")},
			expectedMethod: "POST",
			expectedURL:    "/aaa",
			expectedError:  "request error",
		},
		{
			name:           "RequestWithTLS",
			method:         "put",
			url:            "https://example.com/aaa",
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
			body, err := Request(test.method, test.url, test.data, opts...)

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
