package main

import (
	"crypto/tls"
	"fmt"
	"minivpn/logger"
	"net"
	"server/config"
	"server/handler"
)

func NewServer() (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(config.Cfg.Tls.CertFile, config.Cfg.Tls.KeyFile)
	if err != nil {
		return nil, err
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return tls.Listen("tcp", ":2333", tlsConf)
}

func main() {
	err := config.LoadConfig("config.json")
	if err != nil {
		logger.Log.Fatal("fail to load config:", err)
	}
	listener, err := NewServer()
	if err != nil {
		logger.Log.Fatal(err)
	}
	logger.Log.Info("start listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept fail:", err)
			continue
		}
		go handler.Handle(conn)
	}
	//fmt.Println("I am server")
}
