package ctrago

import (
	"context"
	"testing"
	"time"
)

// 伪造的 Transport 用于单元测试
// 只做基本消息流转，不做真实网络通信

type mockTransport struct {
	sendFn         func(messageType int, data []byte) error
	onMessageFn    func(handler MessageHandler)
	closeFn        func() error
	listenFn       func() error
	setHeartbeatFn func(heartbeatInterval time.Duration, heartbeatFn func() (int, []byte))
}

func (m *mockTransport) Send(messageType int, data []byte) error {
	if m.sendFn != nil {
		return m.sendFn(messageType, data)
	}
	return nil
}
func (m *mockTransport) OnMessage(handler MessageHandler) {
	if m.onMessageFn != nil {
		m.onMessageFn(handler)
	}
}
func (m *mockTransport) Close() error {
	if m.closeFn != nil {
		return m.closeFn()
	}
	return nil
}
func (m *mockTransport) Listen() error {
	if m.listenFn != nil {
		return m.listenFn()
	}
	return nil
}
func (m *mockTransport) SetHeartbeat(heartbeatInterval time.Duration, heartbeatFn func() (int, []byte)) {
	if m.setHeartbeatFn != nil {
		m.setHeartbeatFn(heartbeatInterval, heartbeatFn)
	}
}

func TestClient_SendRequest(t *testing.T) {
	var receivedMsg []byte
	mock := &mockTransport{
		sendFn: func(messageType int, data []byte) error {
			receivedMsg = data
			return nil
		},
	}
	client := NewClientWithTransport(mock, "id", "secret", "token")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 由于 mock，不会有响应，测试超时分支
	_, err := client.SendRequest(ctx, 1, nil)
	if err == nil {
		t.Error("expected timeout error")
	}
	if receivedMsg == nil {
		t.Error("expected Send to be called")
	}
}
