package modules

import (
	"net/http"
)

// ProbeHTTP probes HTTP endpoint
func ProbeHTTP(address string) (bool,error) {
	resp, err := http.Get(address)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, nil
	} 
	return true, nil
}