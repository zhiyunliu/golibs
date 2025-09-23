package xstack

import (
	"bytes"
	"fmt"

	"github.com/zhiyunliu/stack"
)

func GetStack(skip int, opts ...Option) string {
	nopts := make([]stack.Option, 0, len(opts)+1)
	nopts = append(nopts, opts...)
	nopts = append(nopts, stack.WithSkip(skip))
	trace := stack.Trace(nopts...)

	return buildTraceStack(trace)
	//return bytesconv.BytesToString(stack(skip))
}

func buildTraceStack(trace stack.CallStack) string {
	var buffer bytes.Buffer
	for i := range trace {
		buffer.WriteString(fmt.Sprintf("%+n", trace[i]))
		buffer.WriteString("\n\t")
		buffer.WriteString(fmt.Sprintf("%+v", trace[i]))
		buffer.WriteString("\n")
	}

	return buffer.String()
}
