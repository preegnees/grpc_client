package models

// ---> Configuration

type channel struct {
	IdChannel    string `json:"id_channel"`
	AllowedNames string `json:"allowed_names"`
	Send         string `json:"send"`
	Read         string `json:"read"`
}

type Config struct {
	Name     string    `json:"name"`
	Channels []channel `json:"channels"`
}

type IConfig interface {
	Read(path string) (Config, error)
}

// ---> Client

type ClientConf struct {
	Addr              string
	RequestTimeout    int
	KeepaliveInterval int
	Reconnect         bool
	ReconnectTimeout  int
	IdChannel         string
	Name              string
	AllowedNames      string
}

type ICli interface {
	Run(ClientConf) error
	Stop()
}

// ---> Proxy

type TcpProxy struct {
	
}

type ITcpProxy interface {
	Listen(port string) error
}