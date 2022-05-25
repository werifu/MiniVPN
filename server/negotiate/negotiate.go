package negotiate

import (
	"encoding/json"
	"errors"
	"minivpn/logger"
	"minivpn/protocol"
	"net"
	"server/auth"
	"server/config"
	"server/tun"
)

func Negotiate(conn net.Conn) (tun.Tun, error) {
	// auth
	ok := auth.Auth(conn)
	if !ok {
		logger.Log.Error("fail to auth user")
		return tun.Tun{}, errors.New("auth failed")
	}
	logger.Log.Infof("a client auth successfully!")

	// create tun & send ip
	tunnel, err := tun.NewTun()
	if err != nil {
		msg := "fail to create tun: " + err.Error()
		logger.Log.Error(msg)
		pkt, _ := json.Marshal(protocol.HostPacket{
			Ok:   false,
			Host: protocol.Host{},
			Msg:  msg,
		})
		if _, err := conn.Write(pkt); err != nil {
			logger.Log.Error("fail to send fail message to client when shaking failed:", err)
			return tun.Tun{}, err
		}
		return tun.Tun{}, err
	}
	pkt, _ := json.Marshal(protocol.HostPacket{
		Ok: true,
		Host: protocol.Host{
			TunNet:   tunnel.GetClientNet(),
			Intranet: config.Cfg.Intranet,
		},
		Msg: "",
	})
	if _, err := conn.Write(pkt); err != nil {
		logger.Log.Error("fail to write successful handshake packet:", err)
		return tun.Tun{}, err
	}
	return tunnel, nil
}
