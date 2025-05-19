package ctrago

import (
	"context"

	"github.com/yockii/ctrago/openapi"
	"google.golang.org/protobuf/proto"
)

type AccountSymbol struct {
	*Account
}

// AssetsList 资产列表
func (a *AccountSymbol) AssetsList(ctx context.Context, req *openapi.ProtoOAAssetClassListReq) (*openapi.ProtoOAAssetClassListRes, error) {
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_ASSET_CLASS_LIST_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOAAssetClassListRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// SymbolList 获取品种列表
func (a *AccountSymbol) SymbolList(ctx context.Context, req *openapi.ProtoOASymbolsListReq) (*openapi.ProtoOASymbolsListRes, error) {
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_SYMBOLS_LIST_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOASymbolsListRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// SymbolById 根据ID获取品种
func (a *AccountSymbol) SymbolById(ctx context.Context, req *openapi.ProtoOASymbolByIdReq) (*openapi.ProtoOASymbolByIdRes, error) {
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_SYMBOL_BY_ID_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOASymbolByIdRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}

// SymbolsForConversion 获取可兑换品种
func (a *AccountSymbol) SymbolsForConversion(ctx context.Context, req *openapi.ProtoOASymbolsForConversionReq) (*openapi.ProtoOASymbolsForConversionRes, error) {
	respMsg, err := a.client.SendRequest(ctx, uint32(openapi.ProtoOAPayloadType_PROTO_OA_SYMBOLS_FOR_CONVERSION_REQ), req)
	if err != nil {
		return nil, err
	}
	res := &openapi.ProtoOASymbolsForConversionRes{}
	if err := proto.Unmarshal(respMsg.Payload, res); err != nil {
		return nil, err
	}
	return res, nil
}
