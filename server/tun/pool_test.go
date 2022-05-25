package tun

import (
	"bytes"
	"testing"
)

func TestInit(t *testing.T) {
	InitPool(30)
	if len(Pool.tuns) != 64 {
		t.Error("30bit mask net should have a capacity of 64 client")
	}
	if !bytes.Equal(Pool.tuns[0].IPv4[:], []byte{192, 168, 53, 1}) ||
		!bytes.Equal(Pool.tuns[1].IPv4[:], []byte{192, 168, 53, 5}) ||
		!bytes.Equal(Pool.tuns[len(Pool.tuns)-1].IPv4[:], []byte{192, 168, 53, 253}) {
		t.Error("IP generate fail")
	}
}

func TestCap(t *testing.T) {
	InitPool(30)
	for i := 0; i < 64; i++ {
		_, err := Pool.GetTun()
		if err != nil {
			t.Error("cannot get enough tun")
		}
	}
	_, err := Pool.GetTun()
	if err == nil {
		t.Error("should have error")
	}
	err = Pool.ReleaseTun(1)
	if err != nil {
		t.Error("release tun id=1 fail")
	}
	unit, err := Pool.GetTun()
	if err != nil || unit.GetGatewayIP() != "192.168.53.5"{
		t.Error("get released tun fail")
	}
}