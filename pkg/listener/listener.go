package listener

import (
	"fmt"
	"io"
	"net"

	e "streaming/pkg/errors"

	logrus "github.com/sirupsen/logrus"
)

type listener struct {
	log *logrus.Logger
}

func New(logger *logrus.Logger) listener {

	return listener{
		log: logger,
	}
}

func (l *listener) Listen(port string, w io.Writer) error {

	listen, err := net.Listen("tcp", port)
	if err != nil {
		return l.printErr(e.ErrTcpListen, err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return l.printErr(e.ErrAcceptTcpListen, err)
		}
		defer conn.Close()

		buffer := make([]byte, 4096)
		
		_, err = conn.Read(buffer)
		if err != nil {
			return l.printErr(e.ErrReadFromTcpConn, err)
		}
	}
}

func (c *listener) printErr(myerr error, err error) error {

	e := fmt.Errorf("$%v, err:= %v", myerr, err)
	c.log.Error(e)
	return e
}
