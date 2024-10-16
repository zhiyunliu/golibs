package xstack

import "github.com/zhiyunliu/stack"

// Option is a functional option for the stack package.
type Option = stack.Option

// WithDepth sets the maximum depth of the stack trace.
var WithDepth = stack.WithDepth
