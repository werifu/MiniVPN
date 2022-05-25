package config

import "testing"

func TestValidAddr(t *testing.T) {
	validAddrs := []string{"192.156.0.1", "1.1.1.1", "225.225.225.255", "0.1.1.3"}
	for _, validAddr := range validAddrs {
		err := ValidAddr(validAddr)
		if err != nil {
			t.Error(err)
		}
	}

	invalidAddrs := []string{"test", "1.1.1.", "255.256.222.222", "1.4.55.g1"}
	for _, invalidAddr := range invalidAddrs {
		err := ValidAddr(invalidAddr)
		if err == nil {
			t.Error(invalidAddr + " should be invalid address")
		}
	}
}
