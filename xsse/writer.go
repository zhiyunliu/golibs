package xsse

import "io"

type StringWriter io.Writer

// type stringWrapper struct {
// 	io.Writer
// }

// func (w StringWriter) WriteString(str string) (int, error) {
// 	return w.Write([]byte(str))
// }

// func wrapWriter(writer io.Writer) StringWriter {
// 	if w, ok := writer.(StringWriter); ok {
// 		return w
// 	} else {
// 		return stringWrapper{Writer: writer}
// 	}
// }
