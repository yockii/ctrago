package ctrago

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type MessageHandler func(messageType int, data []byte)

type WsClient struct {
	conn     *websocket.Conn
	lock     sync.Mutex
	handlers []MessageHandler

	ticker            *time.Ticker
	heartbeatFn       func() (messageType int, data []byte)
	heartbeatInterval time.Duration
	reconnect         bool
	url               string
	dialer            *websocket.Dialer
	closeCh           chan struct{}
}

// NewWsClientWithHeartbeat 支持心跳和重连的构造方法
func NewWsClientWithHeartbeat(url string, dialer *websocket.Dialer, heartbeatInterval time.Duration, heartbeatFn func() (int, []byte)) *WsClient {
	return &WsClient{
		url:               url,
		dialer:            dialer,
		heartbeatInterval: heartbeatInterval,
		heartbeatFn:       heartbeatFn,
		handlers:          make([]MessageHandler, 0),
		reconnect:         true,
		closeCh:           make(chan struct{}),
	}
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	return &WsClient{
		conn:     conn,
		handlers: make([]MessageHandler, 0),
	}
}

func (c *WsClient) connect() error {
	conn, _, err := c.dialer.Dial(c.url, nil)
	if err != nil {
		return err
	}
	c.lock.Lock()
	c.conn = conn
	c.lock.Unlock()
	return nil
}

func (c *WsClient) Send(messageType int, data []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

func (c *WsClient) OnMessage(handler MessageHandler) {
	c.handlers = append(c.handlers, handler)
}

func (c *WsClient) Listen() error {
	for {
		if c.conn == nil {
			if err := c.connect(); err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
		}
		if c.heartbeatInterval > 0 && c.heartbeatFn != nil {
			if c.ticker == nil {
				c.ticker = time.NewTicker(c.heartbeatInterval)
				go func() {
					for {
						select {
						case <-c.ticker.C:
							msgType, data := c.heartbeatFn()
							c.Send(msgType, data)
						case <-c.closeCh:
							return
						}
					}
				}()
			}
		}
		for {
			messageType, data, err := c.conn.ReadMessage()
			if err != nil {
				c.lock.Lock()
				c.conn.Close()
				c.conn = nil
				c.lock.Unlock()
				if c.reconnect {
					time.Sleep(2 * time.Second)
					break // 跳出内层for，重新连接
				}
				return err
			}
			for _, handler := range c.handlers {
				handler(messageType, data)
			}
		}
	}
}

func (c *WsClient) Close() error {
	close(c.closeCh)
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// WsClient的心跳由Client自动封装时，允许动态设置心跳内容
func (c *WsClient) SetHeartbeat(heartbeatInterval time.Duration, heartbeatFn func() (int, []byte)) {
	c.heartbeatInterval = heartbeatInterval
	c.heartbeatFn = heartbeatFn
	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}
}
