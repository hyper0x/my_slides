package nio

import (
	"net"
	"time"
)

type PubSubContent struct {
	content string
	err     error
}

func (self PubSubContent) Content() string {
	return self.content
}

func (self PubSubContent) Err() error {
	return self.err
}

func NewPubSubContent(content string, err error) PubSubContent {
	return PubSubContent{content: content, err: err}
}

type PubSubListener interface {
	Init(addr string) error
	Listen(handler func(conn net.Conn)) error
	Close() bool
	Addr() net.Addr
}

type PubSubSender interface {
	Init(remoteAddr string, timeout time.Duration) error
	Send(content string) error
	Receive(delim byte) <-chan PubSubContent
	Close() bool
	Addr() net.Addr
	RemoteAddr() net.Addr
}
