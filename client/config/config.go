package config

import (
	"encoding/json"
	"os"
)

type User struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}
type Server struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}
type Tls struct {
	KeyFile  string `json:"key_file"`
	CertFile string `json:"cert_file"`
}
type Config struct {
	User   User   `json:"user"`
	Server Server `json:"server"`
	Tls    Tls    `json:"tls"`
}

var Cfg *Config

func init() {
	Cfg = &Config{}
}

func LoadConfig(path string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &Cfg)
	if err != nil {
		return err
	}
	err = ValidAddr(Cfg.Server.Addr)
	if err != nil {
		return err
	}
	return nil
}
