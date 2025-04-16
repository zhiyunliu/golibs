package engine

type IoWriterWrapper func(bytes []byte) error

func (f IoWriterWrapper) Write(bytes []byte) (int, error) {
	err := f(bytes)
	if err != nil {
		return 0, err
	}
	return len(bytes), nil
}

type ResponseEntity interface {
	StatusCode() int
	Header() map[string]string
	Body() (bytes []byte, err error)
}
