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
func (a *AccountSymbol) AssetsList(ctx context.Context) (*openapi.ProtoOAAssetClassListRes, error) {
	req := &openapi.ProtoOAAssetClassListReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
	}
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
//
// includeArchivedSymbols 是否包含已归档的品种
func (a *AccountSymbol) SymbolList(ctx context.Context, includeArchivedSymbols bool) (*openapi.ProtoOASymbolsListRes, error) {
	req := &openapi.ProtoOASymbolsListReq{
		CtidTraderAccountId:    proto.Int64(a.accountId),
		IncludeArchivedSymbols: proto.Bool(includeArchivedSymbols),
	}
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
//
// symbolIds 需要查询的品种ID列表
// 需要至少传入一个ID
func (a *AccountSymbol) SymbolById(ctx context.Context, symbolIds []int64) (*openapi.ProtoOASymbolByIdRes, error) {
	if len(symbolIds) == 0 {
		return nil, ErrSymbolIdRequired
	}
	req := &openapi.ProtoOASymbolByIdReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		SymbolId:            symbolIds,
	}
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
//
// firstAssetId 第一个资产ID
// lastAssetId 第二个资产ID
// 需要传入两个ID，如EURUSD，firstAssetId=EUR ID，lastAssetId=USD ID
func (a *AccountSymbol) SymbolsForConversion(ctx context.Context, firstAssetId, lastAssetId int64) (*openapi.ProtoOASymbolsForConversionRes, error) {
	if firstAssetId == 0 || lastAssetId == 0 {
		return nil, ErrAssetIdRequired
	}
	req := &openapi.ProtoOASymbolsForConversionReq{
		CtidTraderAccountId: proto.Int64(a.accountId),
		FirstAssetId:        proto.Int64(firstAssetId),
		LastAssetId:         proto.Int64(lastAssetId),
	}
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
