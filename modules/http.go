package modules

import (
	"net/http"
)

func ProbeHTTP(address string) (bool,error) {
	// Make a GET request to the URL
	resp, err := http.Get(address)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else {
		return false, nil
	}
}