package proxy

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	e "streaming/pkg/errors"
)

func (p *Proxy) send(ctx context.Context) error {

	var count int = 0
	for {
		c, err := net.Dial("tcp", fmt.Sprintf(":%s", p.sendPort))
		if err != nil {
			if strings.Contains(err.Error(), "No connection could be made because") {
				time.Sleep(2 * time.Second)
				p.log.Error(fmt.Errorf("Подключение к серверу для отправки данных, попытка:%d", count))
				count++
				continue
			}
			return p.printErr(e.ErrConnToPort, err)
		}
		defer c.Close()
		defer p.reader.Close()
		p.log.Debugf("Подключение к порту %s", p.sendPort)

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				buffer := make([]byte, p.bufferSize)
				defer func() {buffer = nil}()
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
