package errors

import "errors"

var ErrInvalidPathConf = errors.New("invalid path, is not exists")

var ErrReadConfFile = errors.New("Err Read conffile")

var ErrUnmarshalConf = errors.New("Err Unmarshal conffile")

var ErrTcpListen = errors.New("Err listening tcp port")

var ErrAcceptTcpListen = errors.New("Err accepting tcp listen")

var ErrReadFromTcpConn = errors.New("Err reading from tcp conn")

var ErrInvalidBuffer = errors.New("Invalid buffer, mast be < 4096")

var ErrInvalidPorts = errors.New("Invalid ports")

var ErrInvalidName = errors.New("Err invalid name, name is empty")

var ErrInvalidChannels = errors.New("Err invalid channels, channels is empty")

var ErrInvalidIdChannel = errors.New("Err invalid id_channel, id_channel is empty")

var ErrWriteToWriterFromTcpListener = errors.New("Err write to io.writer from tcpListener")

var ErrWriteToWriterFromTcpSender = errors.New("Err write to io.writer from tcpListener")

var ErrConnToPort =  errors.New("Err connect to port")

var ErrReadFromReader = errors.New("Err read from reader")