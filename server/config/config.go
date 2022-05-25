package config

import (
	"encoding/json"
	"os"
)

type AuthUser struct {
	Username   string `json:"username"`
	PasswdHash string `json:"passwd_hash"`
}
type Tls struct {
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}
type Config struct {
	Tls       Tls        `json:"tls"`
	AuthUsers []AuthUser `json:"auth_users"`
	Intranet  string     `json:"intranet"`
}

var Cfg *Config

// AuthUsers is map[username]->passwd_hash
var AuthUsers map[string]string

func init() {
	Cfg = &Config{}
	AuthUsers = make(map[string]string)
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
	for i := 0; i < len(Cfg.AuthUsers); i++ {
		AuthUsers[Cfg.AuthUsers[i].Username] = Cfg.AuthUsers[i].PasswdHash
	}
	return nil
}
