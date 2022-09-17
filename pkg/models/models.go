package models

// ---> Configuration

type Channel struct {
	IdChannel    string `json:"id_channel"`
	AllowedNames string `json:"allowed_names"`
	Ports        string `json:"ports"`
	NewConn      bool   `json:"new"`
}

type Config struct {
	Name      string    `json:"name"`
	Ancillary string    `json:"ancillary"`
	Channels  []Channel `json:"channels"`
	Buffer    int       `json:"buffer"`
}

type IConfig interface {
	Get(path string) (*Config, error)
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

type IListener interface {
	Listen(port string) error
}
