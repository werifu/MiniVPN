package tun

import (
	"errors"
	"fmt"
	"sync"
)

type Unit struct {
	id   int
	IPv4 [4]byte
	mask int
}

type PoolT struct {
	tuns []struct {
		using bool
		Unit
	}
	mut sync.Mutex
}

var Pool PoolT

func init() {
	InitPool([4]byte{192, 168, 53, 0}, 30)
}

func InitPool(net [4]byte, mask int) {
	size := 1 << (8 - (32 - mask))
	offset := 1 << (32 - mask)
	Pool.tuns = make([]struct {
		using bool
		Unit
	}, size)
	for i := 0; i < size; i++ {
		Pool.tuns[i].using = false
		Pool.tuns[i].id = i
		Pool.tuns[i].IPv4 = [4]byte{net[0], net[1], net[2], byte(1 + offset*i)}
		Pool.tuns[i].mask = mask
	}
}

func (p *PoolT) GetTun() (Unit, error) {
	p.mut.Lock()
	defer p.mut.Unlock()
	for i := 0; i < len(p.tuns); i++ {
		if p.tuns[i].using == false {
			p.tuns[i].using = true
			return p.tuns[i].Unit, nil
		}
	}
	return Unit{
		id: -1,
	}, errors.New("no available ip in pool")
}

func (p *PoolT) ReleaseTun(id int) error {
	p.mut.Lock()
	defer p.mut.Unlock()
	if id < 0 || id >= len(p.tuns) {
		return errors.New(fmt.Sprintf("invalid id: %d", id))
	}
	p.tuns[id].using = false
	return nil
}

func (t *Unit) GetGatewayIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", t.IPv4[0], t.IPv4[1], t.IPv4[2], t.IPv4[3])
}

func (t *Unit) GetClientIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", t.IPv4[0], t.IPv4[1], t.IPv4[2], t.IPv4[3]+1)
}

func (t *Unit) GetGatewayNet() string {
	return fmt.Sprintf("%d.%d.%d.%d/%d", t.IPv4[0], t.IPv4[1], t.IPv4[2], t.IPv4[3], t.mask)
}

func (t *Unit) GetClientNet() string {
	return fmt.Sprintf("%d.%d.%d.%d/%d", t.IPv4[0], t.IPv4[1], t.IPv4[2], t.IPv4[3]+1, t.mask)
}
