package xnet

import (
	"net"
	"strconv"
	"strings"
	"time"
)

const StaticLocalIP = "127.0.0.1"

var (
	RetryTimes int           = 10
	Interval   time.Duration = time.Second * 5
)

func getInterfaceAddrs() (addrs []net.Addr, err error) {
	var retryCnt int = 0
	for {
		retryCnt++
		addrs, err = net.InterfaceAddrs()
		if err != nil {
			time.Sleep(Interval)
			continue
		}
		if retryCnt >= RetryTimes {
			break
		}
		if len(addrs) != 0 {
			break
		}
	}
	return
}

// GetLocalIP 获取IP地址
func GetLocalIP(masks ...string) string {
	addrs, _ := getInterfaceAddrs()
	if len(addrs) == 0 {
		return StaticLocalIP
	}

	var ipLst []string
	for _, addr := range addrs {
		ip := addr.(*net.IPNet).IP
		if !IsLocalIPv4(ip) {
			continue
		}
		if len(masks) == 0 {
			return ip.String()
		}
		ipLst = append(ipLst, ip.String())
	}
	for _, ip := range ipLst {
		for _, m := range masks {
			if strings.HasPrefix(ip, m) {
				return ip
			}
		}
	}
	return StaticLocalIP
}

//ipv4: IsLocalIP 检测IP地址是否内网
func IsLocalIPv4Addr(ip string) bool {
	return IsLocalIPv4(net.ParseIP(ip))
}

// ipv4:IsLocalIP 检测IP地址是否内网
func IsLocalIPv4(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

// ExtractHostPort from address
func ExtractHostPort(addr string) (host string, port uint64, err error) {
	var ports string
	host, ports, err = net.SplitHostPort(addr)
	if err != nil {
		return
	}
	port, err = strconv.ParseUint(ports, 10, 16)
	if err != nil {
		return
	}
	return
}

func IsValidIP(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast()
}
