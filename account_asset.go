package ctrago

import (
	"context"

	"github.com/yockii/ctrago/openapi"
	"google.golang.org/protobuf/proto"
)

type AccountAsset struct {
	*Account
}

// AssetList 获取账户资产类别列表
func (a *AccountAsset) AssetList(ctx context.Context) (*openapi.ProtoOAAssetClassListRes, error) {
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
