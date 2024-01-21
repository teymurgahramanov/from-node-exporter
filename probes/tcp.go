package probes

import (
	"log"
	"net"
	"time"
)

func probeTCP(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	defer conn.Close()
	return true
}