package xnet

import (
	"net"
	"strings"
	"time"
)

const StaticLocalMac = "00:1A:2B:3C:4D:5E"

func GetLocalMac(masks ...string) string {
	addrs, _ := GetInterfaces()
	if len(addrs) == 0 {
		return StaticLocalMac
	}
	var first string
	for _, addr := range addrs {
		if len(addr.HardwareAddr) > 0 {
			ads, err := addr.Addrs()
			if err != nil {
				continue
			}
			macV := addr.HardwareAddr.String()
			for _, ad := range ads {
				ip := ad.(*net.IPNet).IP
				if !IsLocalIPv4(ip) {
					continue
				}
				if len(first) <= 0 {
					first = macV
				}
				ipv := ip.String()
				for _, m := range masks {
					if !strings.HasPrefix(ipv, m) {
						continue
					}
					return macV
				}
			}
		}
	}
	if len(first) > 0 {
		return first
	}
	return StaticLocalMac
}

var GetInterfaces = defaultGetInterfaces

func defaultGetInterfaces() (addrs []net.Interface, err error) {
	var retryCnt int = 0
	for {
		retryCnt++
		addrs, err = net.Interfaces()
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
