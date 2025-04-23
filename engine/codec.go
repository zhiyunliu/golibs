package engine

import "net/http"

// IoWriterWrapper is a function that wraps the io.Writer interface.
type IoWriterWrapper func(bytes []byte) error

// Write implements the io.Writer interface.
func (f IoWriterWrapper) Write(bytes []byte) (int, error) {
	err := f(bytes)
	if err != nil {
		return 0, err
	}
	return len(bytes), nil
}

// ResponseEntity is an interface that defines the methods for a response entity.
type ResponseEntity interface {
	StatusCode() int
	Header() map[string]string
	Body() (bytes []byte, err error)
}

// NormalResponse is a struct that implements the ResponseEntity interface.
type NormalResponse struct {
	statusCode int
	header     map[string]string
	body       []byte
}

// NewNormalResponse creates a new NormalResponse instance with the given status code, header, and body.
func NewNormalResponse(statusCode int, header map[string]string, body []byte) ResponseEntity {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	if header == nil {
		header = make(map[string]string)
	}
	if body == nil {
		body = make([]byte, 0)
	}
	return &NormalResponse{
		statusCode: statusCode,
		header:     header,
		body:       body,
	}
}

func (r *NormalResponse) StatusCode() int {
	return r.statusCode
}
func (r *NormalResponse) Header() map[string]string {
	return r.header
}
func (r *NormalResponse) Body() (bytes []byte, err error) {
	return r.body, nil
}
