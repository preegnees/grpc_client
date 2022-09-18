package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	e "streaming/pkg/errors"
	m "streaming/pkg/models"

	logrus "github.com/sirupsen/logrus"
)

type configuration struct {
	log *logrus.Logger
}

func New(logger *logrus.Logger) configuration {

	return configuration{
		log: logger,
	}
}

// добавить потом проверку или пинг сервера

func (c *configuration) Get(path string) (*m.Config, error) {

	conf := m.Config{}

	if _, err := os.Stat(path); err != nil {
		return nil, c.printErr(e.ErrInvalidPathConf, err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, c.printErr(e.ErrReadConfFile, err)
	}

	if err = json.Unmarshal(b, &conf); err != nil {
		return nil, c.printErr(e.ErrUnmarshalConf, err)
	}

	if conf.Buffer > 4096 {
		return nil, c.printErr(e.ErrInvalidBuffer, err)
	}

	if conf.Name == "" {
		return nil, c.printErr(e.ErrInvalidName, err)
	}

	err = c.validPort(conf.Ancillary)
	if err != nil {
		return nil, c.printErr(e.ErrInvalidPorts, err)
	}

	if len(conf.Channels) == 0 {
		return nil, c.printErr(e.ErrInvalidChannels, nil)
	}

	conf.AllPorts = append(conf.AllPorts, conf.Ancillary)

	for _, v := range conf.Channels {
		if len(v.IdChannel) < 8 {
			return nil, c.printErr(e.ErrInvalidIdChannel, nil)
		}
		err = c.validPort(v.Ports)
		if err != nil {
			return nil, c.printErr(e.ErrInvalidPorts, err)
		}
		conf.AllPorts = append(conf.AllPorts, v.Ports)
	}

	return &conf, nil
}

func (c *configuration) printErr(myerr error, err error) error {

	e := fmt.Errorf("$%v, err:= %v", myerr, err)
	c.log.Error(e)
	return e
}

func (c *configuration) validPort(ports string) error {

	var len int = len(strings.Split(ports, " "))
	if len != 2 {
		return fmt.Errorf("Err with ports:%s", ports)
	}

	p1, p2 := strings.Split(ports, " ")[0], strings.Split(ports, " ")[1]
	p1i, err := strconv.Atoi(p1)
	if err != nil {
		return err
	}

	p2i, err := strconv.Atoi(p2)
	if err != nil {
		return err
	}

	if p1i <= 1024 || p1i >= 49151 || p2i <= 1024 || p2i >= 49151 {
		return fmt.Errorf("Err with ports:%s, mast be 1024 < port < 49151", ports)
	}
	return nil
}
