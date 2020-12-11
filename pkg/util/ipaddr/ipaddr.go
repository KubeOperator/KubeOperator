package ipaddr

import (
	"bytes"
	"fmt"
	"github.com/c-robinson/iplib"
	"net"
	"strconv"
)

func GenerateIps(ip string, mask int, startIp string, endIp string) []string {
	var ips []string
	n := iplib.NewNet(net.ParseIP(ip), mask)
	i := n.FirstAddress()
	for {
		if isBiggerThan(i.String(), startIp) >= 0 && isBiggerThan(endIp, i.String()) >= 0 && isAvailableIp(i) {
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

func isAvailableIp(ip net.IP) bool {
	return ip[3] != 0 && ip[3] != 255
}

func ParseMask(num int) (mask string, err error) {
	var buff bytes.Buffer
	for i := 0; i < int(num); i++ {
		buff.WriteString("1")
	}
	for i := num; i < 32; i++ {
		buff.WriteString("0")
	}
	masker := buff.String()
	a, _ := strconv.ParseUint(masker[:8], 2, 64)
	b, _ := strconv.ParseUint(masker[8:16], 2, 64)
	c, _ := strconv.ParseUint(masker[16:24], 2, 64)
	d, _ := strconv.ParseUint(masker[24:32], 2, 64)
	resultMask := fmt.Sprintf("%v.%v.%v.%v", a, b, c, d)
	return resultMask, nil
}
