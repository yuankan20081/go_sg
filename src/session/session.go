// session session.go
package session

import (
	"errors"
	"io"
)

var (
	ErrBufferFull = errors.New("session buffer full")
)

type Session struct {
	sid          int
	conn         io.ReadWriteCloser
	chSendRemote chan []byte
}

func New(conn io.ReadWriteCloser, sid int) *Session {
	return &Session{
		sid:          sid,
		conn:         conn,
		chSendRemote: make(chan []byte, 1024),
	}
}

func (this *Session) Flush() {
	defer this.conn.Close()

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

func (this *Session) Read(p []byte) (n int, err error) {
	return this.conn.Read(p)
}
