package net

import (
	"fmt"
	"net"
	"time"
)

func Ping(ip string) bool {
	_, err := net.DialTimeout("ip4:icmp", ip, time.Duration(1*1000*1000))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
