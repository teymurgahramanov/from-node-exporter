package modules

import (
	"net"
	"time"
)

// ProbeTCP is for probe TCP endpoints
func ProbeTCP(address string, interval int) (bool,error) {
	conn, err := net.DialTimeout("tcp", address, time.Duration(interval)/2*time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return true, nil
}