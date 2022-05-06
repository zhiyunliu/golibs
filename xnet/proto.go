package xnet

import (
	"fmt"
	"strings"
)

//解析协议信息(proto://name)
func Parse(addr string) (proto, name string, err error) {
	addr = strings.TrimSpace(addr)

	parties := strings.SplitN(addr, "://", 2)
	if len(parties) != 2 {
		err = fmt.Errorf("[%s]协议格式错误,正确格式(proto://name)", addr)
		return
	}
	if proto = parties[0]; proto == "" {
		err = fmt.Errorf("[%s]缺少协议proto,正确格式(proto://name)", addr)
		return
	}
	if name = parties[1]; name == "" {
		err = fmt.Errorf("[%s]缺少name,正确格式(proto://name)", addr)
		return
	}
	return
}
