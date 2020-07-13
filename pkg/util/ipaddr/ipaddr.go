package ipaddr

import (
	"github.com/c-robinson/iplib"
	"net"
)

func GenerateIps(ip string, mask int) []string {
	var ips []string
	n := iplib.NewNet(net.ParseIP(ip), mask)
	i := n.FirstAddress()
	for {
		ips = append(ips, i.String())
		i, _ = n.NextIP(i)
		if i.Equal(n.LastAddress()) {
			break
		}
	}
	return ips

}
