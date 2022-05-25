package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	err := LoadConfig("config_test.json")
	if err != nil {
		t.Error(err)
	}
	if Cfg.Tls.KeyFile != "server.key" ||
		Cfg.Tls.CertFile != "server.crt" ||
		len(Cfg.AuthUsers) != 1 ||
		Cfg.AuthUsers[0].Username != "abc" ||
		Cfg.AuthUsers[0].PasswdHash != "114514hash" {
		t.Error("fail to parse")
	}
}
