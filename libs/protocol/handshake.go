package protocol

type HostPacket struct {
	Ok   bool   `json:"ok"`
	Host Host   `json:"host"`
	Msg  string `json:"msg"`
}

type Host struct {
	TunNet   string `json:"tun_net"`
	Intranet string `json:"intranet"`
}

type User struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

type AuthRes struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}
