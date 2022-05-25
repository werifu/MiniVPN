package main

import (
	"client/config"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/songgao/water"
	"io"
	"minivpn/logger"
	"minivpn/protocol"
	"net"
	"os/exec"
	"sync"
)

func tun2sock(iface *water.Interface, conn net.Conn) {
	packet := make([]byte, 2000)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			logger.Log.Fatal(err)
		}
		//logger.Log.Infof("Packet Received from tun, length:%d\n", n)
		n, err = conn.Write(packet[:n])
		if err != nil {
			logger.Log.Error("sock write fail", err)
		}
	}
}
func sock2tun(conn net.Conn, iface *water.Interface) {
	packet := make([]byte, 2000)
	for {
		n, err := conn.Read(packet)
		if err != nil {
			if err == io.EOF {
				logger.Log.Fatal("server closed the connection")
			}
			logger.Log.Fatal(err)
		}
		//logger.Log.Infof("Packet Received from sock, length:%d\n", n)
		n, err = iface.Write(packet[:n])
		if err != nil {
			logger.Log.Error("Iface write fail", err)
		}
	}
}

func connectServer(addr string) (net.Conn, error) {
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		logger.Log.Fatal(err)
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	conn, err := tls.Dial("tcp", addr, tlsConf)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	logger.Log.Info("connected to", conn.RemoteAddr().String())
	return conn, nil
}

func SetupTun(tunNet, intranet string) (*water.Interface, error) {
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	logger.Log.Info("Interface Name:", iface.Name())

	err = exec.Command("ifconfig", "tun0", tunNet, "up").Run()
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	logger.Log.Info(iface.Name() + "'s net: " + tunNet)

	err = exec.Command("route", "add", "-net", intranet, iface.Name()).Run()
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	logger.Log.Info("intranet:", intranet)

	return iface, nil
}

func Handshake(conn net.Conn) (protocol.HostPacket, error) {
	buf := make([]byte, 1024)

	pkt := protocol.HostPacket{}
	n, err := conn.Read(buf)
	if err != nil {
		logger.Log.Error("fail to receive handshake packet from server:", err)
		return pkt, err
	}

	err = json.Unmarshal(buf[:n], &pkt)
	if err != nil {
		logger.Log.Error("fail to parse handshake packet", buf[:n], err)
		return pkt, err
	}
	if !pkt.Ok {
		logger.Log.Error("fail to handshake, message:", pkt.Msg)
		return pkt, errors.New(pkt.Msg)
	}
	return pkt, nil
}

func BiForward(conn net.Conn, iface *water.Interface) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		sock2tun(conn, iface)
		wg.Done()
	}()
	go func() {
		tun2sock(iface, conn)
		wg.Done()
	}()
	wg.Wait()
}

func Auth(conn net.Conn, username, passwd string) error {
	buf, _ := json.Marshal(protocol.User{
		Username: username,
		Passwd:   passwd,
	})
	_, err := conn.Write(buf)
	if err != nil {
		logger.Log.Error("fail to write auth information:", err)
		return err
	}

	res := make([]byte, 1024)
	authRes := protocol.AuthRes{}
	n, err := conn.Read(res)
	if err != nil {
		logger.Log.Error("fail to read auth response:", err)
		return err
	}
	err = json.Unmarshal(res[:n], &authRes)
	if err != nil {
		logger.Log.Error("fail to parse auth res:", err)
		return err
	}
	if authRes.Ok == true {
		logger.Log.Info("succeed to auth!")
		return nil
	} else {
		logger.Log.Error("fail to auth:", authRes.Msg)
		return errors.New(authRes.Msg)
	}
}

func main() {
	err := config.LoadConfig("config.json")
	if err != nil {
		logger.Log.Fatal(err)
	}

	// connect
	addr := fmt.Sprintf("%s:%d", config.Cfg.Server.Addr, config.Cfg.Server.Port)
	conn, err := connectServer(addr)
	if err != nil {
		logger.Log.Fatal("fail to connect to server:", err)
	}

	// auth
	username := config.Cfg.User.Username
	passwd := config.Cfg.User.Passwd
	err = Auth(conn, username, passwd)
	if err != nil {
		logger.Log.Fatal(err)
	}

	// get ip
	pkt, err := Handshake(conn)
	if err != nil {
		logger.Log.Fatal(err)
	}

	// setup tunnel
	iface, err := SetupTun(pkt.Host.TunNet, pkt.Host.Intranet)
	if err != nil {
		logger.Log.Fatal(err)
	}

	// handler
	BiForward(conn, iface)
}
