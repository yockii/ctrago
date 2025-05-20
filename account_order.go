package ctrago

import (
	"context"

	"github.com/yockii/ctrago/openapi"
	"google.golang.org/protobuf/proto"
)

type AccountOrder struct {
	*Account
}

// NewOrder 下单
func (a *AccountOrder) NewOrder(ctx context.Context, symbolId int64, orderType openapi.ProtoOAOrderType, tradeSide openapi.ProtoOATradeSide, volume int64, orderOption *OrderOption) (*openapi.ProtoOAExecutionEvent, error) {
	if symbolId <= 0 {
		return nil, ErrSymbolIdRequired
	}
	if volume <= 0 {
		return nil, ErrVolumeRequired
	}
	req := &openapi.ProtoOANewOrderReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		SymbolId:            proto.Int64(symbolId),
		OrderType:           &orderType,
		TradeSide:           &tradeSide,
		Volume:              proto.Int64(volume),
	}
	if orderOption != nil {
		if orderOption.limitPrice > 0 {
			req.LimitPrice = proto.Float64(orderOption.limitPrice)
		}
		if orderOption.stopPrice > 0 {
			req.StopPrice = proto.Float64(orderOption.stopPrice)
		}
		if orderOption.timeInForce != 0 {
			req.TimeInForce = &orderOption.timeInForce
		}
		if orderOption.expirationTimestamp > 0 {
			req.ExpirationTimestamp = proto.Int64(orderOption.expirationTimestamp)
		}
		if orderOption.stopLoss > 0 {
			req.StopLoss = proto.Float64(orderOption.stopLoss)
		}
		if orderOption.takeProfit > 0 {
			req.TakeProfit = proto.Float64(orderOption.takeProfit)
		}
		if orderOption.comment != "" {
			req.Comment = proto.String(orderOption.comment)
		}
		if orderOption.baseSlippagePrice > 0 {
			req.BaseSlippagePrice = proto.Float64(orderOption.baseSlippagePrice)
		}
		if orderOption.label != "" {
			req.Label = proto.String(orderOption.label)
		}
		if orderOption.positionId > 0 {
			req.PositionId = proto.Int64(orderOption.positionId)
		}
		if orderOption.clientOrderId != "" {
			req.ClientOrderId = proto.String(orderOption.clientOrderId)
		}
		if orderOption.relativeStopLoss > 0 {
			req.RelativeStopLoss = proto.Int64(orderOption.relativeStopLoss)
		}
		if orderOption.relativeTakeProfit > 0 {
			req.RelativeTakeProfit = proto.Int64(orderOption.relativeTakeProfit)
		}
		if orderOption.guaranteedStopLoss {
			req.GuaranteedStopLoss = proto.Bool(orderOption.guaranteedStopLoss)
		}
		if orderOption.trailingStopLoss {
			req.TrailingStopLoss = proto.Bool(orderOption.trailingStopLoss)
		}
		if orderOption.stopTriggerMethod != 0 {
			req.StopTriggerMethod = &orderOption.stopTriggerMethod
		}
	}
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_NEW_ORDER_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAExecutionEvent{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// CancelOrder 撤单
func (a *AccountOrder) CancelOrder(ctx context.Context, orderId int64) (*openapi.ProtoOAExecutionEvent, error) {
	req := &openapi.ProtoOACancelOrderReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		OrderId:             proto.Int64(orderId),
	}
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_CANCEL_ORDER_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAExecutionEvent{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AmendOrder 修改订单
func (a *AccountOrder) AmendOrder(ctx context.Context, orderId int64, orderOption *AmendOrderOption) (*openapi.ProtoOAExecutionEvent, error) {
	req := &openapi.ProtoOAAmendOrderReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		OrderId:             proto.Int64(orderId),
	}
	if orderOption != nil {
		if orderOption.volume > 0 {
			req.Volume = proto.Int64(orderOption.volume)
		}
		if orderOption.limitPrice > 0 {
			req.LimitPrice = proto.Float64(orderOption.limitPrice)
		}
		if orderOption.stopPrice > 0 {
			req.StopPrice = proto.Float64(orderOption.stopPrice)
		}
		if orderOption.expirationTimestamp > 0 {
			req.ExpirationTimestamp = proto.Int64(orderOption.expirationTimestamp)
		}
		if orderOption.stopLoss > 0 {
			req.StopLoss = proto.Float64(orderOption.stopLoss)
		}
		if orderOption.takeProfit > 0 {
			req.TakeProfit = proto.Float64(orderOption.takeProfit)
		}
		if orderOption.slippageInPoints > 0 {
			req.SlippageInPoints = proto.Int32(orderOption.slippageInPoints)
		}
		if orderOption.relativeStopLoss > 0 {
			req.RelativeStopLoss = proto.Int64(orderOption.relativeStopLoss)
		}
		if orderOption.relativeTakeProfit > 0 {
			req.RelativeTakeProfit = proto.Int64(orderOption.relativeTakeProfit)
		}
		if orderOption.guaranteedStopLoss {
			req.GuaranteedStopLoss = proto.Bool(orderOption.guaranteedStopLoss)
		}
		if orderOption.trailingStopLoss {
			req.TrailingStopLoss = proto.Bool(orderOption.trailingStopLoss)
		}
		if orderOption.stopTriggerMethod != 0 {
			req.StopTriggerMethod = &orderOption.stopTriggerMethod
		}
	}

	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_AMEND_ORDER_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAExecutionEvent{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AmendOrderPositionSlip 修改订单止损止盈
func (a *AccountOrder) AmendOrderPositionSltp(ctx context.Context, positionId int64, amendPositionSLTPOption *AmendPositionSLTPOption) (*openapi.ProtoOAExecutionEvent, error) {
	if positionId <= 0 {
		return nil, ErrPositionIdRequired
	}
	req := &openapi.ProtoOAAmendPositionSLTPReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		PositionId:          proto.Int64(positionId),
	}
	if amendPositionSLTPOption != nil {
		if amendPositionSLTPOption.stopLoss > 0 {
			req.StopLoss = proto.Float64(amendPositionSLTPOption.stopLoss)
		}
		if amendPositionSLTPOption.takeProfit > 0 {
			req.TakeProfit = proto.Float64(amendPositionSLTPOption.takeProfit)
		}
		if amendPositionSLTPOption.guaranteedStopLoss {
			req.GuaranteedStopLoss = proto.Bool(amendPositionSLTPOption.guaranteedStopLoss)
		}
		if amendPositionSLTPOption.trailingStopLoss {
			req.TrailingStopLoss = proto.Bool(amendPositionSLTPOption.trailingStopLoss)
		}
		if amendPositionSLTPOption.stopLossTriggerMethod != 0 {
			req.StopLossTriggerMethod = &amendPositionSLTPOption.stopLossTriggerMethod
		}
	}

	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_AMEND_POSITION_SLTP_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAExecutionEvent{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ClosePosition 平仓
func (a *AccountOrder) ClosePosition(ctx context.Context, positionId, volume int64) (*openapi.ProtoOAExecutionEvent, error) {
	if positionId <= 0 {
		return nil, ErrPositionIdRequired
	}
	if volume <= 0 {
		return nil, ErrVolumeRequired
	}
	req := &openapi.ProtoOAClosePositionReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		PositionId:          proto.Int64(positionId),
		Volume:              proto.Int64(volume),
	}
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_CLOSE_POSITION_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAExecutionEvent{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}
