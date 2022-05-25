package tun

import (
	"errors"
	"github.com/songgao/water"
	"minivpn/logger"
	"os/exec"
)

type Tun struct {
	Iface 	*water.Interface
	Unit
}

func NewTun() (tun Tun, err error) {
	// get a net in pool
	unit, err := Pool.GetTun()
	if err != nil {
		return Tun{}, err
	}

	// create tun interface
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		logger.Log.Error("create tun fail", err)
		_ = Pool.ReleaseTun(unit.id)
		return Tun{}, err
	}
	net := unit.GetGatewayNet()
	toExec := "ifconfig " + iface.Name() + " " + net + " up"

	// setup tun's ip
	err = exec.Command("ifconfig", iface.Name(), net, "up").Run()
	if err != nil {
		_ = Pool.ReleaseTun(unit.id)
		_ = iface.Close()
		return Tun{}, err
	}
	logger.Log.Infof("Create Tun Interface Name: %s\n", iface.Name())
	logger.Log.Info("Exec:", toExec)
	return Tun{
		Iface: iface,
		Unit:  unit,
	}, nil
}

func (t *Tun) Close() error {
	err1 := t.Iface.Close()
	err2 := Pool.ReleaseTun(t.id)
	if err1 != nil {
		logger.Log.Error("fail to close tun interface", err1)
	}
	if err2 != nil {
		logger.Log.Error("fail to release tun ip", err2)
	}
	if err1 != nil || err2 != nil {
		return errors.New(err1.Error() + "|" + err2.Error())
	}
	return nil
}