package ctrago

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yockii/ctrago/openapi"
	"google.golang.org/protobuf/proto"
)

// Transport 通信抽象接口
// Send: 发送消息
// OnMessage: 注册消息回调
// Close: 关闭连接
// Listen: 启动消息循环（如有需要）
type Transport interface {
	Send(messageType int, data []byte) error
	OnMessage(handler MessageHandler)
	Close() error
	Listen() error
	SetHeartbeat(heartbeatInterval time.Duration, heartbeatFn func() (int, []byte))
}

type ResponseHandler func(*openapi.ProtoMessage)

// 修改Client结构体，底层通信改为Transport接口
// 并将NewClient的第一个参数类型由*WsClient改为Transport
type Client struct {
	transport     Transport
	msgId         uint64
	lock          sync.Mutex
	pending       map[string]chan *openapi.ProtoMessage
	eventHandlers map[uint32][]ResponseHandler

	clientId     string
	clientSecret string
	accessToken  string
}

func NewClientWithTransport(transport Transport, clientId, clientSecret, accessToken string) *Client {
	c := &Client{
		transport:     transport,
		pending:       make(map[string]chan *openapi.ProtoMessage),
		eventHandlers: make(map[uint32][]ResponseHandler),
		clientId:      clientId,
		clientSecret:  clientSecret,
		accessToken:   accessToken,
	}
	transport.OnMessage(c.handleMessage)
	return c
}

// NewClientWithWebsocket 使用 WebSocket 创建 Client
func NewClientWithWebsocket(addr, clientId, clientSecret, accessToken string, heartbeatInterval time.Duration) *Client {
	ws := NewWsClientWithHeartbeat(addr, nil, 0, nil)
	client := NewClientWithTransport(ws, clientId, clientSecret, accessToken)
	ws.SetHeartbeat(heartbeatInterval, func() (int, []byte) {
		hb := &openapi.ProtoHeartbeatEvent{}
		data, _ := proto.Marshal(hb)
		return websocket.BinaryMessage, data
	})
	go ws.Listen()
	return client
}

// NewClientWithTcp 使用 TCP 创建 Client
func NewClientWithTcp(addr, clientId, clientSecret, accessToken string) (*Client, error) {
	tcp, err := NewTcpClient(addr)
	if err != nil {
		return nil, err
	}
	client := NewClientWithTransport(tcp, clientId, clientSecret, accessToken)
	go tcp.Listen()
	return client, nil
}

// NewClient 默认使用 WebSocket 方式创建 Client
func NewClient(addr, clientId, clientSecret, accessToken string, heartbeatInterval time.Duration) *Client {
	return NewClientWithWebsocket(addr, clientId, clientSecret, accessToken, heartbeatInterval)
}

func (c *Client) nextMsgId() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.msgId++
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), c.msgId)
}

func (c *Client) SendRequest(ctx context.Context, payloadType uint32, payload proto.Message) (*openapi.ProtoMessage, error) {
	msgId := c.nextMsgId()
	data, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}
	msg := &openapi.ProtoMessage{
		PayloadType: proto.Uint32(payloadType),
		Payload:     data,
		ClientMsgId: proto.String(msgId),
	}
	raw, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	ch := make(chan *openapi.ProtoMessage, 1)
	c.lock.Lock()
	c.pending[msgId] = ch
	c.lock.Unlock()
	err = c.transport.Send(websocket.BinaryMessage, raw)
	if err != nil {
		return nil, err
	}
	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) handleMessage(messageType int, data []byte) {
	if messageType != websocket.BinaryMessage {
		return
	}
	msg := &openapi.ProtoMessage{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return
	}
	if msg.ClientMsgId != nil && *msg.ClientMsgId != "" {
		c.lock.Lock()
		ch, ok := c.pending[*msg.ClientMsgId]
		if ok {
			delete(c.pending, *msg.ClientMsgId)
		}
		c.lock.Unlock()
		if ok {
			ch <- msg
		}
		return
	}
	// 事件推送
	if msg.PayloadType != nil {
		c.lock.Lock()
		handlers := c.eventHandlers[*msg.PayloadType]
		c.lock.Unlock()
		for _, h := range handlers {
			h(msg)
		}
	}
}

func (c *Client) OnEvent(payloadType uint32, handler ResponseHandler) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.eventHandlers[payloadType] = append(c.eventHandlers[payloadType], handler)
}

func (c *Client) Close() error {
	return c.transport.Close()
}

// ApplicationAuth 应用鉴权
func (c *Client) ApplicationAuth(ctx context.Context) (*openapi.ProtoOAApplicationAuthRes, error) {
	req := &openapi.ProtoOAApplicationAuthReq{
		ClientId:     &c.clientId,
		ClientSecret: &c.clientSecret,
	}
	respMsg, err := c.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAApplicationAuthRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Version 获取OpenAPI版本
func (c *Client) Version(ctx context.Context) (*openapi.ProtoOAVersionRes, error) {
	req := &openapi.ProtoOAVersionReq{}
	respMsg, err := c.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_VERSION_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAVersionRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetAccountList 获取账户列表
func (c *Client) GetAccountList(ctx context.Context) (*openapi.ProtoOAGetAccountListByAccessTokenRes, error) {
	req := &openapi.ProtoOAGetAccountListByAccessTokenReq{
		AccessToken: &c.accessToken,
	}
	respMsg, err := c.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_GET_ACCOUNTS_BY_ACCESS_TOKEN_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAGetAccountListByAccessTokenRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// RefreshToken 刷新token
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*openapi.ProtoOARefreshTokenRes, error) {
	req := &openapi.ProtoOARefreshTokenReq{
		RefreshToken: &refreshToken,
	}
	respMsg, err := c.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_REFRESH_TOKEN_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOARefreshTokenRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Account 返回账户操作对象
func (c *Client) Account(accountId int64) *Account {
	return &Account{
		client:    c,
		accountId: accountId,
	}
}
