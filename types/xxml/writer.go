package xxml

import (
	"bytes"
	"fmt"
	"io"

	"github.com/zhiyunliu/golibs/types"
)

//Writer Writer
type Writer struct {
	buffer   *bytes.Buffer
	eleStack *types.Stack
}

//NewWriter NewWriter
func NewWriter() *Writer {
	return &Writer{
		buffer:   bytes.NewBuffer(nil),
		eleStack: types.NewStack(),
	}
}

//WriteStart WriteStart
func (w *Writer) WriteStart(ele string) *Writer {
	w.buffer.WriteString(fmt.Sprintf("<%s>", ele))
	w.eleStack.Push(ele)
	return w
}

//WriteEnd WriteEnd
func (w *Writer) WriteEnd() *Writer {
	ele := w.eleStack.Pop()
	if ele == nil {
		panic(fmt.Errorf("must after begin"))
	}
	w.buffer.WriteString(fmt.Sprintf("</%s>", ele))
	return w
}

//WriteValue WriteValue
func (w *Writer) WriteValue(val string) *Writer {
	w.buffer.WriteString(val)
	return w
}

func (w *Writer) String() string {
	return w.buffer.String()
}

func (w *Writer) WriteTo(wto io.Writer) (n int, err error) {
	return wto.Write(w.buffer.Bytes())
}
