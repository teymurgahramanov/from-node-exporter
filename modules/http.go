package modules

import (
	"net/http"
)

func ProbeHTTP(address string) (bool,error) {
	resp, err := http.Get(address)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else {
		return false, nil
	}
}