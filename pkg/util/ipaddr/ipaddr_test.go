package ipaddr

import (
	"fmt"
	"testing"
)

func TestGenerateIps(t *testing.T) {
	ips := GenerateIps("172.16.10.0", 24)
	fmt.Println(len(ips))
}
