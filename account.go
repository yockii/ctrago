package ctrago

import (
	"context"

	"github.com/yockii/ctrago/openapi"
	"google.golang.org/protobuf/proto"
)

type Account struct {
	client    *Client
	accountId int64
}

func NewAccount(client *Client, accountId int64) *Account {
	return &Account{
		client:    client,
		accountId: accountId,
	}
}

// 示例：账户登录
func (a *Account) Auth(ctx context.Context) (*openapi.ProtoOAAccountAuthRes, error) {
	req := &openapi.ProtoOAAccountAuthReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		AccessToken:         proto.String(a.client.accessToken),
	}
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAAccountAuthRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// 便于聚合各类账户操作
func (a *Account) Order() *AccountOrder {
	return &AccountOrder{Account: a}
}

func (a *Account) Symbol() *AccountSymbol {
	return &AccountSymbol{Account: a}
}

func (a *Account) Trader() *AccountTrader {
	return &AccountTrader{Account: a}
}

func (a *Account) Asset() *AccountAsset {
	return &AccountAsset{Account: a}
}
