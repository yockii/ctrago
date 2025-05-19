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
func (a *AccountOrder) NewOrder(ctx context.Context, req *openapi.ProtoOANewOrderReq) (*openapi.ProtoOAExecutionEvent, error) {
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
func (a *AccountOrder) CancelOrder(ctx context.Context, req *openapi.ProtoOACancelOrderReq) (*openapi.ProtoOAExecutionEvent, error) {
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
func (a *AccountOrder) AmendOrder(ctx context.Context, req *openapi.ProtoOAAmendOrderReq) (*openapi.ProtoOAExecutionEvent, error) {
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
func (a *AccountOrder) AmendOrderPositionSlip(ctx context.Context, req *openapi.ProtoOAAmendPositionSLTPReq) (*openapi.ProtoOAExecutionEvent, error) {
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
func (a *AccountOrder) ClosePosition(ctx context.Context, req *openapi.ProtoOAClosePositionReq) (*openapi.ProtoOAExecutionEvent, error) {
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
