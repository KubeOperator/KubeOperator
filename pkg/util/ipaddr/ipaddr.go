package ipaddr

import (
	"bytes"
	"github.com/c-robinson/iplib"
	"net"
)

func GenerateIps(ip string, mask int, startIp string, endIp string) []string {
	var ips []string
	n := iplib.NewNet(net.ParseIP(ip), mask)
	i := n.FirstAddress()
	for {
		if isBiggerThan(i.String(), startIp) >= 0 && isBiggerThan(endIp, i.String()) >= 0 {
			ips = append(ips, i.String())
		}
		i, _ = n.NextIP(i)
		if i.Equal(n.LastAddress()) {
			break
		}
	}
	return ips
}

func isBiggerThan(a string, b string) int {
	aIp := net.ParseIP(a)
	bIp := net.ParseIP(b)
	return bytes.Compare(aIp, bIp)
}
