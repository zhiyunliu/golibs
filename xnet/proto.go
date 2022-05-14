package xnet

import (
	"fmt"
	"net/url"
	"strings"
)

//解析协议信息(proto://name)
func Parse(addr string) (proto, name string, err error) {
	addr = strings.TrimSpace(addr)

	val, err := url.Parse(addr)
	if err != nil {
		err = fmt.Errorf("[%s]协议格式错误,正确格式(proto://name)", addr)
		return
	}
	proto = val.Scheme
	name = val.Host
	return
}
