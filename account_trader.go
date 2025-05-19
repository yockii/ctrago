package ctrago

import (
	"context"

	"github.com/yockii/ctrago/openapi"
	"google.golang.org/protobuf/proto"
)

type AccountTrader struct {
	*Account
}

// Trader 获取账户信息
func (a *AccountTrader) Trader(ctx context.Context, req *openapi.ProtoOATraderReq) (*openapi.ProtoOATraderRes, error) {
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_TRADER_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOATraderRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// 你可以继续扩展更多账户信息相关方法
