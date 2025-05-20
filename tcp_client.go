package ctrago

import (
	"net"
	"sync"
	"time"
)

// TcpClient 实现 Transport 接口，支持 cTrader OpenAPI 的 TCP 通信
// 仅实现最基础的 Send/OnMessage/Listen/Close，心跳可选实现

type TcpClient struct {
	conn     net.Conn
	lock     sync.Mutex
	handlers []MessageHandler
	closeCh  chan struct{}
}

func NewTcpClient(addr string) (*TcpClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &TcpClient{
		conn:     conn,
		handlers: make([]MessageHandler, 0),
		closeCh:  make(chan struct{}),
	}, nil
}

func (c *TcpClient) Send(messageType int, data []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.conn.Write(data)
	return err
}

func (c *TcpClient) OnMessage(handler MessageHandler) {
	c.handlers = append(c.handlers, handler)
}

func (c *TcpClient) Listen() error {
	buf := make([]byte, 4096)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			return err
		}
		for _, handler := range c.handlers {
			handler(0, buf[:n])
		}
	}
}

func (c *TcpClient) Close() error {
	close(c.closeCh)
	return c.conn.Close()
}

func (c *TcpClient) SetHeartbeat(heartbeatInterval time.Duration, heartbeatFn func() (int, []byte)) {
	// TCP心跳可选实现，暂留空
}

var _ Transport = (*TcpClient)(nil)
