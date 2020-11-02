package logagent

import (
	"errors"
	// "fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	logf "github.com/alonegrowing/purple/pkg/basic/kernel/logging"
)

const (
	recvChanLength int = 2000
	sendChanLength int = 2000
)

var ClientDisconnectedError = errors.New("Client has already closed")
var ClientWriteBlockingError = errors.New("write packet was blocking")

type Conn struct {
	client     net.Conn
	remoteAddr string

	closeOnce sync.Once

	closeChan chan bool
	closed    int32

	// packetRecvChan chan []byte
	dataSendChan chan []byte
	// logg         *logf.Logger
}

func NewLogConn(remoteAddr string) (*Conn, error) {

	conn, err := net.Dial("udp", remoteAddr)
	if err != nil {
		// fmt.Println("连接服务端失败:", err.Error())
		// logf.Error("logagent conn server error,", remoteAddr)
		return nil, err
	}
	logConn := &Conn{
		client:     conn,
		remoteAddr: remoteAddr,
		closeChan:  make(chan bool),
		closeOnce:  sync.Once{},
		// packetRecvChan: make(chan []byte, recvChanLength),
		dataSendChan: make(chan []byte, sendChanLength),
	}

	go logConn.handleLoop()
	// go logConn.readLoop()
	// go logConn.writeLoop()

	// logf.Info("logagent conn server succ,", remoteAddr)
	// fmt.Println("logagent conn server succ,", remoteAddr)
	return logConn, nil
}

func (conn *Conn) handleLoop() {
	defer func() {
		recover()
		conn.Close()
		// logf.Error("logagent handlerLoop error,", conn.remoteAddr)
	}()

	for {
		select {
		case <-conn.closeChan:
			return
			// case <-conn.packetRecvChan:

		}
	}
}

func (conn *Conn) writeLoop() {

	defer func() {
		recover()
		conn.Close()
		logf.Error("logagent writeLoop error,", conn.remoteAddr)
	}()

	for {

		select {

		case <-conn.closeChan:
			return

		case p, ok := <-conn.dataSendChan:
			if !ok {
				return
			}
			if conn.Closed() {
				return
			}

			// fmt.Println("write:", len(p))
			if _, err := conn.client.Write(p); err != nil {
				// logf.Error("logagent writeLoop error0,", err)
				// fmt.Println("logagent writeLoop error0,", err)
				conn.Close()
				return
			}
		}
	}
}

func (conn *Conn) Closed() bool {
	return atomic.LoadInt32(&conn.closed) == 1
}

func (conn *Conn) Write(b []byte, timeout time.Duration) error {

	if _, err := conn.client.Write(b); err != nil {
		conn.Close()
		return err
	}

	// var err error

	// err = nil
	// if conn.Closed() {
	// 	err = ClientDisconnectedError
	// 	return err
	// }
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		err = ClientDisconnectedError
	// 	}
	// }()

	// if timeout == 0 {

	// 	select {
	// 	case <-conn.closeChan:
	// 		err = ClientDisconnectedError
	// 		return err

	// 	case conn.dataSendChan <- b:
	// 		return err
	// 	default:
	// 		err = ClientWriteBlockingError
	// 		return err

	// 	}

	// } else {
	// 	select {
	// 	case <-conn.closeChan:
	// 		err = ClientDisconnectedError
	// 		return err
	// 	case conn.dataSendChan <- b:
	// 		//fmt.Println("wtire datasendchan succ")
	// 		return err
	// 	case <-time.After(timeout):
	// 		err = ClientWriteBlockingError
	// 		return err
	// 	}
	// }

	return nil
}

func (conn *Conn) Close() {
	conn.closeOnce.Do(func() {
		atomic.StoreInt32(&conn.closed, 1)
		conn.client.Close()
		close(conn.closeChan)

		close(conn.dataSendChan)
		// close(conn.packetRecvChan)
		// logf.Error("logagent close ,", conn.remoteAddr)

	})
}
