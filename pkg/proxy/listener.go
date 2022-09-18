package proxy

import (
	"context"
	"fmt"
	"net"

	e "streaming/pkg/errors"
)

func (p *Proxy) listen(ctx context.Context) error {

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", p.recvPort))
	if err != nil {
		return p.printErr(e.ErrTcpListen, err)
	}
	defer listen.Close()
	defer p.writer.Close()
	p.log.Debugf("Прослушивание порта %s", p.recvPort)

	conn, err := listen.Accept()
	if err != nil {
		return p.printErr(e.ErrAcceptTcpListen, err)
	}
	defer conn.Close()
	p.log.Debugf("Подключение к порту %s", p.recvPort)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			buffer := make([]byte, p.bufferSize)
			defer func() {buffer = nil}()
			_, err = conn.Read(buffer)
			if err != nil {
				return p.printErr(e.ErrReadFromTcpConn, err)
			}
			_, err = p.writer.Write(buffer)
			if err != nil {
				return p.printErr(e.ErrWriteToWriterFromTcpListener, err)
			}
		}
	}
}
