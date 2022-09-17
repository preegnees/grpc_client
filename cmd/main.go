package main

import (
	"fmt"
	"log"
	"os"

	c "streaming/pkg/configuration"
	m "streaming/pkg/models"
	ltcp "streaming/pkg/listener"

	logrus "github.com/sirupsen/logrus"
)

func main() {
	
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	conf, err := readConf(logrus.New())
	if err != nil {
		return err
	}
	log.Println(conf)


	


	err = listen(conf, logrus.New())
	if err != nil {
		return err
	}

	return nil
}

func listen(conf *m.Config, log *logrus.Logger) error {
	
	listener := ltcp.New(log)

	return nil
}

func readConf(log *logrus.Logger) (*m.Config, error) {
	if len(os.Args[1:]) != 1 {
		return nil, fmt.Errorf("Args is empty")
	}
	
	configuration := c.New(log)
	conf, err := configuration.Get(os.Args[1:][0])
	if err != nil {
		log.Fatal(err)
	}

	return conf, nil
}