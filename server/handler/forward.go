package handler

import (
	"github.com/songgao/water"
	"io"
	"minivpn/logger"
	"net"
	"server/negotiate"
	"sync"
)

const MaxClientRetry = 3

func Handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()
		logger.Log.Info("a connection closed")
	}()

	tunnel, err := negotiate.Negotiate(conn)
	if err != nil {
		logger.Log.Error("fail to negotiate with client", err)
		return
	}

	// handler
	biForward(conn, tunnel.Iface)

	// finish
	err = tunnel.Close()
	if err != nil {
		logger.Log.Error("fail to close tunnel", err)
	}
}

func biForward(conn net.Conn, iface *water.Interface) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	closeSig := make(chan bool)
	go func() {
		sock2tun(conn, iface, closeSig)
		wg.Done()
	}()
	go func() {
		tun2sock(iface, conn, closeSig)
		wg.Done()
	}()
	wg.Wait()
}

func closeSigSafely(closeSig chan bool) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	close(closeSig)
}

func tun2sock(iface *water.Interface, conn net.Conn, closeSig chan bool) {
	defer closeSigSafely(closeSig)
	packet := make([]byte, 2000)
	retryTimes := 0
	for retryTimes <= MaxClientRetry {
		select {
		case <-closeSig:
			// ch has been closed
			return
		default:
			n, err := iface.Read(packet)
			if err != nil {
				logger.Log.Error(err)
				retryTimes++
				continue
			}
			//logger.Log.Info("packet received from tun:", packet[:n])
			n, err = conn.Write(packet[:n])
			if err != nil {
				logger.Log.Error("sock write fail", err)
				retryTimes++
				continue
			}
			retryTimes = 0
		}
	}
}

func sock2tun(conn net.Conn, iface *water.Interface, closeSig chan bool) {
	defer closeSigSafely(closeSig)
	packet := make([]byte, 2000)
	retryTimes := 0
	for retryTimes <= MaxClientRetry {
		select {
		case <-closeSig:
			// ch has been closed
			return
		default:
			n, err := conn.Read(packet)
			if err != nil {
				logger.Log.Error(err)
				if err == io.EOF {
					return
				}
				retryTimes++
				continue
			}
			//logger.Log.Info("packet received from sock:", packet[:n])
			n, err = iface.Write(packet[:n])
			if err != nil {
				logger.Log.Error("Iface write fail", err)
				retryTimes++
				continue
			}
			retryTimes = 0
		}
	}
}
