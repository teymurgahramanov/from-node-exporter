package modules

import (
	"net"
)

func ProbeICMP(address string) (bool,error) {
	ipAddr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return false, err
	}

	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	msg := []byte("hello")
	_, err = conn.Write(msg)
	if err != nil {
		return false, err
	}

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		return false, err
	}

	return true, nil
}