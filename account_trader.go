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
func (a *AccountTrader) Trader(ctx context.Context) (*openapi.ProtoOATraderRes, error) {
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_TRADER_REQ), &openapi.ProtoOATraderReq{
		CtidTraderAccountId: &a.accountId,
	})
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOATraderRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Reconcile 获取账户当前持仓和挂单
//
// returnProtectionOrders 是否返回保护单
func (a *AccountTrader) Reconcile(ctx context.Context, returnProtectionOrders bool) (*openapi.ProtoOAReconcileRes, error) {
	// 设置账户ID
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_RECONCILE_REQ), &openapi.ProtoOAReconcileReq{
		CtidTraderAccountId:    proto.Int64(a.accountId),
		ReturnProtectionOrders: proto.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAReconcileRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// DealList 获取账户历史成交明细
func (a *AccountTrader) DealList(ctx context.Context, fromTimestamp, toTimestamp int64, maxRows int32) (*openapi.ProtoOADealListRes, error) {
	req := &openapi.ProtoOADealListReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
	}

	if fromTimestamp >= 0 {
		req.FromTimestamp = proto.Int64(fromTimestamp)
	}
	if toTimestamp >= 0 && toTimestamp > fromTimestamp {
		req.ToTimestamp = proto.Int64(toTimestamp)
	}
	if maxRows > 0 {
		req.MaxRows = proto.Int32(maxRows)
	}
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_DEAL_LIST_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOADealListRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// OrderList 获取账户历史订单明细
func (a *AccountTrader) OrderList(ctx context.Context, fromTimestamp, toTimestamp int64) (*openapi.ProtoOAOrderListRes, error) {
	req := &openapi.ProtoOAOrderListReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
	}
	if fromTimestamp >= 0 {
		req.FromTimestamp = proto.Int64(fromTimestamp)
	}
	if toTimestamp >= 0 && toTimestamp > fromTimestamp {
		req.ToTimestamp = proto.Int64(toTimestamp)
	}

	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_ORDER_LIST_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAOrderListRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ExpectedMargin 获取账户保证金预估
func (a *AccountTrader) ExpectedMargin(ctx context.Context, symbolId int64, volumes []int64) (*openapi.ProtoOAExpectedMarginRes, error) {
	if symbolId == 0 {
		return nil, ErrSymbolIdRequired
	}
	req := &openapi.ProtoOAExpectedMarginReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		SymbolId:            proto.Int64(symbolId),
	}
	if len(volumes) > 0 {
		req.Volume = volumes
	}

	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_EXPECTED_MARGIN_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAExpectedMarginRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// CashFlowHistoryList 获取账户资金流水（充值/提现历史）
func (a *AccountTrader) CashFlowHistoryList(ctx context.Context, fromTimestamp, toTimestamp int64) (*openapi.ProtoOACashFlowHistoryListRes, error) {
	if fromTimestamp <= 0 {
		return nil, ErrFromTimestampRequired
	}
	if toTimestamp <= 0 {
		return nil, ErrToTimestampRequired
	}
	if toTimestamp <= fromTimestamp || toTimestamp-fromTimestamp > 604800000 {
		return nil, ErrTimestampRange
	}
	req := &openapi.ProtoOACashFlowHistoryListReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
	}
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_CASH_FLOW_HISTORY_LIST_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOACashFlowHistoryListRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// 你可以继续扩展更多账户信息相关方法
