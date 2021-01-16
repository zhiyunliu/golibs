package xxml

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/zhiyunliu/golibs/types"
)

//Writer Writer
type Writer struct {
	buffer        *bytes.Buffer
	eleStack      *types.Stack
	compress      bool
	lastwriteType writeType
}
type writeType int

const writeTypeStart writeType = 1
const writeTypeVal writeType = 2
const writeTypeEnd writeType = 3

//NewWriter NewWriter
func NewWriter(compress bool) *Writer {
	return &Writer{
		compress: compress,
		buffer:   bytes.NewBuffer(nil),
		eleStack: types.NewStack(),
	}
}

//WriteStart WriteStart
func (w *Writer) WriteStart(ele string) *Writer {
	w.writetab()
	ele = escapeChars(ele)
	w.buffer.WriteString(fmt.Sprintf("<%s>", ele))
	w.eleStack.Push(ele)
	w.lastwriteType = writeTypeStart
	return w
}

//WriteEnd WriteEnd
func (w *Writer) WriteEnd() *Writer {
	ele := w.eleStack.Pop()
	if ele == nil {
		panic(fmt.Errorf("must after begin"))
	}
	if !w.compress && w.lastwriteType == writeTypeEnd {
		w.buffer.WriteString("\r")
		w.buffer.WriteString(strings.Repeat("\t", w.eleStack.Len()))
	}
	w.buffer.WriteString(fmt.Sprintf("</%s>", ele))
	w.lastwriteType = writeTypeEnd
	return w
}

//WriteValue WriteValue
func (w *Writer) WriteValue(val string) *Writer {
	val = escapeChars(val)
	w.buffer.WriteString(val)
	w.lastwriteType = writeTypeVal
	return w
}

func (w *Writer) String() string {
	return w.buffer.String()
}

// WriteTo  WriteTo
func (w *Writer) WriteTo(wto io.Writer) (n int, err error) {
	return wto.Write(w.buffer.Bytes())
}

func (w *Writer) writetab() {
	if w.compress {
		return
	}
	if w.eleStack.Len() > 0 {
		w.buffer.WriteString("\r")
	}
	w.buffer.WriteString(strings.Repeat("\t", w.eleStack.Len()))
}
