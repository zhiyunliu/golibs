package xnet

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/zhiyunliu/golibs/xlog"
)

type PortPicker func(logger xlog.Logger, localIp string, begin, end, step int64) int64

var (
	rnd      = rand.New(rand.NewSource(time.Now().Unix()))
	portPick = map[string]PortPicker{
		"rand": randPortPicker,
		"seq":  seqPortPicker,
	}
)

// :8080
// 2000-30000
// 2000-30000:rand
// 2000-30000:seq:2
func GetAvaliableAddr(logger xlog.Logger, localIp string, addr string) (newAddr string, err error) {
	//没有指定范围
	if !strings.Contains(addr, "-") {
		return addr, nil
	}
	method := "rand"
	step := int64(1)
	parties := strings.SplitN(addr, ":", 3)
	if len(parties) == 2 {
		addr = parties[0]
		method = parties[1]
	}
	if len(parties) == 3 {
		addr = parties[0]
		method = parties[1]
		step, _ = strconv.ParseInt(parties[2], 10, 32)
	}

	parties = strings.SplitN(addr, "-", 2)
	begin, err := strconv.ParseInt(parties[0], 10, 32)
	if err != nil || begin <= 0 {
		return "", fmt.Errorf("指定端口配置错误:%s", addr)
	}
	end, err := strconv.ParseInt(parties[1], 10, 32)
	if err != nil || end <= 0 || end < begin {
		return "", fmt.Errorf("指定端口配置错误:%s", addr)
	}
	call, ok := portPick[method]
	if !ok {
		return "", fmt.Errorf("指定端口配置错误:%s", addr)
	}
	np := call(logger, localIp, begin, end, step)
	if np == 0 {
		return "", fmt.Errorf("未获取到有效的端口,请检查配置:%s", addr)
	}
	newAddr = fmt.Sprintf("%s:%d", localIp, np)
	return
}

func randPortPicker(logger xlog.Logger, localIp string, begin, end, step int64) int64 {
	for {
		np := rnd.Int63n(end-begin) + begin
		logger.Logf(xlog.LevelInfo, "检测端口(rand):%d", np)
		if !ScanPort("TCP", localIp, np) {
			return np
		}
	}
}

func seqPortPicker(logger xlog.Logger, localIp string, begin, end, step int64) int64 {
	for np := begin; np < end; np++ {
		logger.Logf(xlog.LevelInfo, "检测端口(seq):%d", np)
		if !ScanPort("TCP", localIp, np) {
			return np
		}
	}
	return 0
}

func ScanPort(protocol string, hostname string, port int64) bool {
	p := strconv.FormatInt(port, 10)
	addr := net.JoinHostPort(hostname, p)
	conn, err := net.DialTimeout(protocol, addr, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
