package xnet

import (
	"fmt"
	"math/rand"
	"net"
	"regexp"
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

	_portRegexpPattern1 = regexp.MustCompile(`^:(\d+)$`)
	_portRegexpPattern2 = regexp.MustCompile(`^\[(\d+),(\d+)\)(\:rand)?$`)
	_portRegexpPattern3 = regexp.MustCompile(`^\[(\d+),(\d+)\)\:(seq|rand)\:(\d+)$`)

	_portRegexpList []*addrRegexp = []*addrRegexp{
		{
			regexp: _portRegexpPattern1,
			buildAddr: func(logger xlog.Logger, localIp string, subList []string) (string, error) {
				return fmt.Sprintf("%s:%s", localIp, subList[0]), nil
			},
		},
		{
			regexp: _portRegexpPattern2,
			buildAddr: func(logger xlog.Logger, localIp string, subList []string) (string, error) {
				begin, err := strconv.ParseInt(subList[0], 10, 32)
				if err != nil {
					return "", err
				}
				end, err := strconv.ParseInt(subList[1], 10, 32)
				if err != nil {
					return "", err
				}
				nport := randPortPicker(logger, localIp, begin, end, 1)
				return fmt.Sprintf("%s:%d", localIp, nport), nil
			},
		},
		{
			regexp: _portRegexpPattern3,
			buildAddr: func(logger xlog.Logger, localIp string, subList []string) (string, error) {
				begin, err := strconv.ParseInt(subList[0], 10, 32)
				if err != nil {
					return "", err
				}
				end, err := strconv.ParseInt(subList[1], 10, 32)
				if err != nil {
					return "", err
				}
				pickmethod := subList[2]
				pickmethod = strings.Trim(pickmethod, ":")

				step, err := strconv.ParseInt(subList[3], 10, 32)
				if err != nil {
					return "", err
				}
				call, ok := portPick[pickmethod]
				if !ok {
					return "", fmt.Errorf("指定端口配置错误")
				}
				nport := call(logger, localIp, begin, end, step)
				if nport == 0 {
					return "", fmt.Errorf("未获取到有效的端口,请检查配置")
				}
				return fmt.Sprintf("%s:%d", localIp, nport), nil
			},
		},
	}
)

// :8080
// [2000,30000)
// [2000,30000):rand
// [2000,30000):seq:5
func GetAvaliableAddr(logger xlog.Logger, localIp string, addr string) (newAddr string, err error) {
	for _, reg := range _portRegexpList {
		if !reg.MatchString(addr) {
			continue
		}
		return reg.BuildAddr(logger, localIp, addr)
	}
	return "", fmt.Errorf("指定端口配置错误:%s", addr)
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

type addrRegexp struct {
	regexp    *regexp.Regexp
	buildAddr func(logger xlog.Logger, localIp string, subList []string) (string, error)
}

func (p addrRegexp) MatchString(addr string) bool {
	return p.regexp.MatchString(addr)
}
func (p addrRegexp) BuildAddr(logger xlog.Logger, localIp string, addr string) (newAddr string, err error) {
	sublist := p.regexp.FindStringSubmatch(addr)
	newAddr, err = p.buildAddr(logger, localIp, sublist[1:])
	if err != nil {
		return newAddr, fmt.Errorf("%w:%s", err, addr)
	}
	return newAddr, nil
}
