package models

import "io"

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
	Server    string    `json:"server"`
	AllPorts  []string
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
	Reader            *io.PipeReader
	Writer            *io.PipeWriter
}

type ICli interface {
	Run(ClientConf) error
	Stop()
}
