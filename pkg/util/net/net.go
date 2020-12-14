package net

import (
	"crypto/tls"
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


func TcpPing(addr string,isTls bool) error  {
	if isTls{
		conn,err:=tls.DialWithDialer(&net.Dialer{Timeout: 5*time.Second},"tcp",addr,&tls.Config{InsecureSkipVerify: true})
		if err!=nil{
			return err
		}
		defer conn.Close()
	}else {
		conn,err:=net.DialTimeout("tcp",addr,5*time.Second)
		if err!=nil{
			return err
		}
		defer conn.Close()
	}
	return nil
}