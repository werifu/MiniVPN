package auth

import (
	"testing"
)

func TestCheckUser(t *testing.T) {
	hash := Hash("114514", "abc")
	//fmt.Println(hash)
	if hash != "8632207854a19adc9fa641d8b824d8cedebfb151c79e9fe1c2b0db1af182fed4" {
		t.Error("hash not equal")
	}
}
