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

type ResponseHandler func(*openapi.ProtoMessage)

type Client struct {
	ws            *WsClient
	msgId         uint64
	lock          sync.Mutex
	pending       map[string]chan *openapi.ProtoMessage
	eventHandlers map[uint32][]ResponseHandler

	clientId     string
	clientSecret string
	accessToken  string
}

func NewClient(ws *WsClient, clientId, clientSecret, accessToken string) *Client {
	c := &Client{
		ws:            ws,
		pending:       make(map[string]chan *openapi.ProtoMessage),
		eventHandlers: make(map[uint32][]ResponseHandler),
		clientId:      clientId,
		clientSecret:  clientSecret,
		accessToken:   accessToken,
	}
	ws.OnMessage(c.handleMessage)
	return c
}

// 新的构造函数，自动封装ProtoHeartbeatEvent心跳
func NewClientWithHeartbeat(url, clientId, clientSecret, accessToken string, heartbeatInterval time.Duration) *Client {
	ws := NewWsClientWithHeartbeat(url, nil, 0, nil) // 先不设置心跳
	client := NewClient(ws, clientId, clientSecret, accessToken)
	ws.SetHeartbeat(heartbeatInterval, func() (int, []byte) {
		hb := &openapi.ProtoHeartbeatEvent{}
		data, _ := proto.Marshal(hb)
		return websocket.BinaryMessage, data
	})
	go ws.Listen()
	return client
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
	err = c.ws.Send(websocket.BinaryMessage, raw)
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
	return c.ws.Close()
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
