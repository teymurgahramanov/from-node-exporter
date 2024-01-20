package probe

import (
	"regexp"
)

func isValidIP(address string) bool {
	ipHostPattern := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}+$`)
	return ipHostPattern.MatchString(address)
}

func isValidIPPort(address string) bool {
	ipPortPattern := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+$`)
	return ipPortPattern.MatchString(address)
}

func isValidHostPort(address string) bool {
	hostPortPattern := regexp.MustCompile(`^[a-zA-Z0-9.-]+:\d+$`)
	return hostPortPattern.MatchString(address)
}