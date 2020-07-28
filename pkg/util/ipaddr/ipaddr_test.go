package ipaddr

import (
	"fmt"
	"testing"
)

func TestGenerateIps(t *testing.T) {
	ips := GenerateIps("172.16.10.0", 24, "172.16.10.2", "172.16.10.23")
	for _, ip := range ips {
		fmt.Println(ip)
	}
}
