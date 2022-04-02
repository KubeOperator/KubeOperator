package net

import (
	"net"
	"time"
)

func TcpPing(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	defer func() {
		_ = conn.Close()
	}()
	return err
}
