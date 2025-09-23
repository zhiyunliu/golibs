package engine

type ResponseWriter interface {
	Status() int
	Size() int
	Written() bool
	WriteHeader(code int)
	Header() Header
	Write(p []byte) (n int, err error)
	WriteString(string) (int, error)
	Flush() error
}
