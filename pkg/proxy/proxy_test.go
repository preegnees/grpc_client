package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	logrus "github.com/sirupsen/logrus"
)

func Test_ConnectToServerIfHeIsAvailable(t *testing.T) {
	pr1, pw1 := io.Pipe()
	pr3, pw3 := io.Pipe()

	p := New(logrus.New(), "44441", "44442", pr1, nil, 4096)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)

	go func() {
		err := runTcpServer(ctx, "44441", pw3, 4096)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err := p.send(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := pr3.Read(buf)
			if err != nil {
				fmt.Println(err)
				cancel()
				return
			}
			fmt.Printf("Колличество байт:%d, message:%s\n", n, buf)
		}
	}()

	tiker := time.NewTicker(1 * time.Second)

	count := 0

	for {
	select {
	case <-tiker.C:
		n, err := pw1.Write([]byte(fmt.Sprintf("%v\n", time.Now())))
		fmt.Println(n)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("круг %d\n", count)
		count++
	case err := <-errCh:
		fmt.Println(err)
		cancel()
	}
	}
}

func runTcpServer(ctx context.Context, port string, w *io.PipeWriter, buffer int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer listen.Close()
	defer w.Close()

	fmt.Printf("Запустился:%v\n", listen)

	conn, err := listen.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Подключился клиент, %v\n", conn)

	for {
		buffer := make([]byte, buffer)
		_, err = conn.Read(buffer)
		if err != nil {
			return err
		}

		_, err = w.Write(buffer)
		if err != nil {
			return err
		}
	}
}
