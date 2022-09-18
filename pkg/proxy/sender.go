package proxy

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	e "streaming/pkg/errors"
)

var dialErr = "No connection could be made because the target machine actively refused it"

var count int = 0

func (p *Proxy) send(ctx context.Context) error {

	for {
		c, err := net.Dial("tcp", fmt.Sprintf(":%s", p.sendPort))
		if err != nil {
			if strings.Contains(err.Error(), dialErr) {
				time.Sleep(2 * time.Second)
				p.log.Error(fmt.Errorf("Подключение к серверу для отправки данных, попытка:%d", count))
				count++
				continue
			}
			return p.printErr(e.ErrConnToPort, err)
		}

		defer c.Close()
		defer p.reader.Close()

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				buffer := make([]byte, p.bufferSize)
				_, err := p.reader.Read(buffer)
				if err != nil {
					return p.printErr(e.ErrReadFromReader, err)
				}
				_, err = c.Write(buffer)
				if err != nil {
					return p.printErr(e.ErrWriteToWriterFromTcpSender, err)
				}
			}
		}
	}
}
