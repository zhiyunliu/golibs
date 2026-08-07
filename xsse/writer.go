package xsse

import "io"

type stringWriter interface {
	io.Writer
	WriteString(string) (int, error)
}

type stringWrapper struct {
	io.Writer
}

func (w stringWrapper) WriteString(str string) (int, error) {
	return w.Write([]byte(str))
}

func wrapWriter(writer io.Writer) stringWriter {
	if w, ok := writer.(stringWriter); ok {
		return w
	} else {
		return stringWrapper{Writer: writer}
	}
}
