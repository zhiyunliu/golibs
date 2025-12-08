package xstack

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

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

func PkgFilePath(frame *runtime.Frame) string {
	pre := pkgPrefix(frame.Function)
	post := pathSuffix(frame.File)
	if pre == "" {
		return post
	}
	return pre + "/" + post
}

// pkgPrefix returns the import path of the function's package with the final
// segment removed.
func pkgPrefix(funcName string) string {
	const pathSep = "/"
	end := strings.LastIndex(funcName, pathSep)
	if end == -1 {
		return ""
	}
	return funcName[:end]
}

// pathSuffix returns the last two segments of path.
func pathSuffix(path string) string {
	const pathSep = "/"
	lastSep := strings.LastIndex(path, pathSep)
	if lastSep == -1 {
		return path
	}
	return path[strings.LastIndex(path[:lastSep], pathSep)+1:]
}
