package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"minivpn/logger"
	"minivpn/protocol"
	"net"
	"server/config"
)

func Auth(conn net.Conn) bool {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		logger.Log.Error("fail to read auth information from client:", err)
		res, _ := json.Marshal(protocol.AuthRes{
			Ok:  false,
			Msg: "fail to read auth information",
		})
		_, _ = conn.Write(res)
		return false
	}
	user := &protocol.User{}
	err = json.Unmarshal(buf[:n], user)
	if err != nil {
		logger.Log.Error("fail to parse auth json:", err)
		res, _ := json.Marshal(protocol.AuthRes{
			Ok:  false,
			Msg: "fail to parse auth",
		})
		_, _ = conn.Write(res)
		return false
	}
	ok := CheckUser(user.Username, user.Passwd)
	if !ok {
		logger.Log.Error("fail to check user")
		res, _ := json.Marshal(protocol.AuthRes{
			Ok:  false,
			Msg: "fail to check user",
		})
		_, _ = conn.Write(res)
	}
	res, _ := json.Marshal(protocol.AuthRes{
		Ok:  true,
		Msg: "",
	})
	_, _ = conn.Write(res)
	return ok
}

func CheckUser(username string, passwd string) bool {
	if passwdHashInSvr, ok := config.AuthUsers[username]; ok {
		hashResult := Hash(passwd, username)
		if hashResult == passwdHashInSvr {
			return true
		}
		return false
	} else {
		return false
	}
}

// Hash returns passwdHash (not so safe but easy)
func Hash(passwd string, salt string) string {
	hash := sha256.Sum256([]byte(passwd + salt))
	return fmt.Sprintf("%x", hash[:])
}
