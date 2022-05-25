package config

import (
	"errors"
	"regexp"
)

func ValidAddr(addr string) error {
	// copy from https://ihateregex.io/expr/ip/
	matched, err := regexp.Match("(\\b25[0-5]|\\b2[0-4]\\d|\\b[01]?\\d\\d?)(\\.(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)){3}", []byte(addr))
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid ipv4 address:" + addr)
	}
	return nil
}
