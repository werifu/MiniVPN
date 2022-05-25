module server

go 1.18

replace minivpn/logger => ../libs/logger

require minivpn/logger v0.0.1

replace minivpn/protocol => ../libs/protocol

require (
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8
	minivpn/protocol v0.0.1
)

require (
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/withmandala/go-log v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/sys v0.0.0-20220519141025-dcacdad47464 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
)
