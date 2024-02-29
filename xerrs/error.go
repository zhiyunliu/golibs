package xerrs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xtypes"
)

type Xerror interface {
	GetCode() int
	GetData() xtypes.XMap
	Error() string
}

type xErr struct {
	Code   int
	Data   xtypes.XMap
	Ignore bool
	error  error
}

// GetCode 获取错误码
func (a xErr) GetCode() int {
	return a.Code
}

// GetCode 获取错误码
func (a xErr) GetData() xtypes.XMap {
	return a.Data
}

// GetError 获取错误信息
func (a xErr) GetError() error {
	return a
}

// GetError 获取错误信息
func (a xErr) String() string {
	bytes, _ := json.Marshal(a)
	return bytesconv.BytesToString(bytes)
}

// GetCode 获取错误码
func (a xErr) Error() string {
	return a.error.Error()
}

func (a xErr) Is(e error) bool {
	return errors.Is(a.error, e)
}
func (a xErr) As(target interface{}) bool {
	return errors.As(a.error, target)
}

func New(err error, opts ...Option) Xerror {
	return NewCode(GetCode(err, 901), err, opts...)
}

// Newf 创建错误对象
func Newf(f string, args ...interface{}) Xerror {
	return New(fmt.Errorf(f, args...))
}

func NewCode(code int, err error, opts ...Option) Xerror {
	xe := &xErr{
		Code:  code,
		error: err,
	}
	for i := range opts {
		opts[i](xe)
	}
	return xe
}

func GetCode(err error, def ...int) int {
	switch v := err.(type) {
	case Xerror:
		return v.GetCode()
	default:
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
}

func GetData(err error) xtypes.XMap {
	switch v := err.(type) {
	case Xerror:
		return v.GetData()
	default:
		return xtypes.XMap{}
	}
}
