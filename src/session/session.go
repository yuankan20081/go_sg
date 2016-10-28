// session session.go
package session

import (
	"errors"
	"io"
	"time"
)

var (
	ErrBufferFull = errors.New("session buffer full")
)

type OnReadHandler interface {
	OnRead([]byte, int)
}
type OnReadHandleFunc func([]byte, int)

func (fn OnReadHandleFunc) OnRead(p []byte, sid int) {
	fn(p, sid)
}

type Session struct {
	sid          int
	conn         io.ReadWriteCloser
	chSendRemote chan []byte
	readHandler  OnReadHandler
}

func New(conn io.ReadWriteCloser, sid int) *Session {
	return &Session{
		sid:          sid,
		conn:         conn,
		chSendRemote: make(chan []byte, 1024),
	}
}

func (this *Session) OnRead(handler OnReadHandler) {
	this.readHandler = handler
}
func (this *Session) OnReadFunc(fn func([]byte, int)) {
	this.readHandler = OnReadHandleFunc(fn)
}

func (this *Session) Serve(timeout time.Duration) {
	defer this.conn.Close()

	go this.parseRecv(this.conn)

	for preSendMessage := range this.chSendRemote {
		this.conn.Write(preSendMessage)
	}
}

func (this *Session) Stop() {
	close(this.chSendRemote)
}

func (this *Session) Write(p []byte) (n int, err error) {
	select {
	case this.chSendRemote <- p:
	default:
		return 0, ErrBufferFull
	}
	return len(p), nil
}

func (this *Session) parseRecv(r io.Reader) {
	var buf [4096]byte
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				this.Stop()
			}
			break
		}
		if this.readHandler != nil {
			this.readHandler.OnRead(buf[:n], this.sid)
		}
	}
}
