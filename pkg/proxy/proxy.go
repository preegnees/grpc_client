package proxy

import (
	"context"
	"fmt"
	"io"

	logrus "github.com/sirupsen/logrus"
)

type Proxy struct {
	log        *logrus.Logger
	sendPort   string
	recvPort   string
	reader     *io.PipeReader
	writer     *io.PipeWriter
	bufferSize int
}

func New(logger *logrus.Logger, sPort, rPort string, r *io.PipeReader, w *io.PipeWriter, bufferSize int) Proxy {

	return Proxy{
		log:        logger,
		sendPort:   sPort,
		recvPort:   rPort,
		reader:     r,
		writer:     w,
		bufferSize: bufferSize,
	}
}

func (p *Proxy) Run(ctx context.Context, errCh chan<- error) {

	p.log.Debugf("Запущена функция Run() с sendPort=%s, recvPort=%s", p.sendPort, p.recvPort)
	go func(ctx context.Context, eCh chan<- error) {
		err := p.listen(ctx)
		if err != nil {
			if eCh != nil {
				eCh <- err
				return
			}
			return
		}
		return
	}(ctx, errCh)

	go func(ctx context.Context, eCh chan<- error) {
		err := p.send(ctx)
		if err != nil {
			if eCh != nil {
				eCh <- err
				return
			}
			return
		}
		return
	}(ctx, errCh)
}

func (p *Proxy) printErr(myerr error, err error) error {

	e := fmt.Errorf("$%v, err:= %v", myerr, err)
	p.log.Error(e)
	return e
}
