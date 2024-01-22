package modules

import (
	"net"
	"time"
)

func ProbeTCP(address string) (bool,error) {
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return false, nil
	}
	defer conn.Close()
	return true, nil
}