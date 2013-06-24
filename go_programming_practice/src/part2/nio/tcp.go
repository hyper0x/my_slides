package nio

import (
	"errors"
	"net"
	"sync"
	"time"
)

type TcpListener struct {
	listener net.Listener
	active   bool
	lock     *sync.Mutex
}

func (self *TcpListener) Init(addr string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		return nil
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	self.listener = ln
	self.active = true
	return nil
}

func (self *TcpListener) Listen(handler func(conn net.Conn)) error {
	if !self.active {
		return errors.New("Send Error: Uninitialized listener!")
	}
	go func() {
		for {
			conn, err := self.listener.Accept()
			if err != nil {
				Logger().Errorf("Listener: Accept Request Error: %s\n", err)
				continue
			}
			go handler(conn)
		}
	}()
	return nil
}

func (self *TcpListener) Close() bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		self.listener.Close()
		self.active = false
		return true
	} else {
		return false
	}
}

func (self *TcpListener) Addr() net.Addr {
	if self.active {
		return self.listener.Addr()
	} else {
		return nil
	}
}

func NewTcpListener() PubSubListener {
	return &TcpListener{lock: new(sync.Mutex)}
}

type TcpSender struct {
	active bool
	lock   *sync.Mutex
	conn   net.Conn
}

func (self *TcpSender) Init(remoteAddr string, timeout time.Duration) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.active {
		conn, err := net.DialTimeout("tcp", remoteAddr, timeout)
		if err != nil {
			return err
		}
		self.conn = conn
		self.active = true
	}
	return nil
}

func (self *TcpSender) Send(content string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.active {
		return errors.New("Send Error: Uninitialized sender!")
	}
	_, err := WriteToTcp(self.conn, content)
	return err
}

func (self *TcpSender) Receive(delim byte) <-chan PubSubContent {
	respChan := make(chan PubSubContent, 1)
	go func(conn net.Conn, ch chan<- PubSubContent) {
		content, err := ReadFromTcp(conn, DELIM)
		ch <- NewPubSubContent(content, err)
	}(self.conn, respChan)
	return respChan
}

func (self *TcpSender) Addr() net.Addr {
	if self.active {
		return self.conn.LocalAddr()
	} else {
		return nil
	}
}

func (self *TcpSender) RemoteAddr() net.Addr {
	if self.active {
		return self.conn.RemoteAddr()
	} else {
		return nil
	}
}

func (self *TcpSender) Close() bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		self.conn.Close()
		self.active = false
		return true
	} else {
		return false
	}
}

func NewTcpSender() PubSubSender {
	return &TcpSender{lock: new(sync.Mutex)}
}
