package xsse

import "io"

type StringWriter interface {
	io.Writer
	WriteString(string) (int, error)
}

type stringWrapper struct {
	io.Writer
}

func (w stringWrapper) WriteString(str string) (int, error) {
	return w.Write([]byte(str))
}

func wrapWriter(writer io.Writer) StringWriter {
	if w, ok := writer.(StringWriter); ok {
		return w
	} else {
		return stringWrapper{Writer: writer}
	}
}
