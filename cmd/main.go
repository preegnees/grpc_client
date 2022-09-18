package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"os/signal"
    "syscall"

	cli "streaming/pkg/cli"
	cnf "streaming/pkg/configuration"
	m "streaming/pkg/models"
	p "streaming/pkg/proxy"

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

	err = runProxy(*conf, logrus.New())
	if err != nil {
		return err
	}

	return nil
}

func readConf(log *logrus.Logger) (*m.Config, error) {
	if len(os.Args[1:]) != 1 {
		return nil, fmt.Errorf("Args is empty")
	}

	configuration := cnf.New(log)
	conf, err := configuration.Get(os.Args[1:][0])
	if err != nil {
		log.Fatal(err)
	}

	return conf, nil
}

func runProxy(conf m.Config, log *logrus.Logger) error {

	proxies := make([]p.Proxy, 0)

	errCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	// служебный канал потом реализовать

	for _, v := range conf.Channels {
		prToTcpSender, pwFromRecvRpc := io.Pipe()
		prToRpcSender, pwFromRecvTcp := io.Pipe()
		portRead := strings.Split(v.Ports, " ")[0]
		portWrite := strings.Split(v.Ports, " ")[1]

		proxy := p.New(
			log,
			portWrite,
			portRead,
			prToTcpSender,
			pwFromRecvTcp,
			conf.Buffer,
		)

		proxies = append(proxies, proxy)

		proxy.Run(ctx, errCh)

		go func(v m.Channel, eCh chan error) {
			c := cli.New()
			err := c.Run(m.ClientConf{
				Addr:              conf.Server,
				RequestTimeout:    15,
				KeepaliveInterval: 600,
				Reconnect:         true,
				ReconnectTimeout:  15,
				IdChannel:         v.IdChannel,
				Name:              conf.Name,
				AllowedNames:      fmt.Sprintf("%s,%s", v.AllowedNames, conf.Name),
				Reader:            prToRpcSender,
				Writer:            pwFromRecvRpc,
			})
			if err != nil {
				if eCh != nil {
					eCh <- err
					return
				}
				return
			}
			return
		}(v, errCh)
	}

	c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        os.Exit(1)
    }()

	select {
	case err := <-errCh:
		close(errCh)
		errCh = nil
		cancel()
		return err
	}
}
