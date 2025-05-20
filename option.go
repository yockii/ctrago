package ctrago

import "github.com/yockii/ctrago/openapi"

type BaseOrderOption struct {
	limitPrice          float64
	stopPrice           float64
	expirationTimestamp int64
	stopLoss            float64
	takeProfit          float64
	comment             string
	label               string
	positionId          int64
	clientOrderId       string
	relativeStopLoss    int64
	relativeTakeProfit  int64
	guaranteedStopLoss  bool
	trailingStopLoss    bool
	stopTriggerMethod   openapi.ProtoOAOrderTriggerMethod
}

type AmendOrderOption struct {
	BaseOrderOption
	volume           int64
	slippageInPoints int32
}

func (o *AmendOrderOption) WithVolume(volume int64) *AmendOrderOption {
	o.volume = volume
	return o
}
func (o *AmendOrderOption) WithSlippageInPoints(points int32) *AmendOrderOption {
	o.slippageInPoints = points
	return o
}

//////////////////////////////////////
type OrderOption struct {
	BaseOrderOption
	timeInForce       openapi.ProtoOATimeInForce
	baseSlippagePrice float64
}

func (o *OrderOption) WithTimeInForce(tif openapi.ProtoOATimeInForce) *OrderOption {
	o.timeInForce = tif
	return o
}
func (o *OrderOption) WithBaseSlippagePrice(price float64) *OrderOption {
	o.baseSlippagePrice = price
	return o
}

////////////////////
func (o *BaseOrderOption) WithLimitPrice(price float64) *BaseOrderOption {
	o.limitPrice = price
	return o
}
func (o *BaseOrderOption) WithStopPrice(price float64) *BaseOrderOption {
	o.stopPrice = price
	return o
}
func (o *BaseOrderOption) WithExpirationTimestamp(ts int64) *BaseOrderOption {
	o.expirationTimestamp = ts
	return o
}
func (o *BaseOrderOption) WithStopLoss(price float64) *BaseOrderOption {
	o.stopLoss = price
	return o
}
func (o *BaseOrderOption) WithTakeProfit(price float64) *BaseOrderOption {
	o.takeProfit = price
	return o
}
func (o *BaseOrderOption) WithComment(comment string) *BaseOrderOption {
	o.comment = comment
	return o
}
func (o *BaseOrderOption) WithLabel(label string) *BaseOrderOption {
	o.label = label
	return o
}
func (o *BaseOrderOption) WithPositionId(id int64) *BaseOrderOption {
	o.positionId = id
	return o
}
func (o *BaseOrderOption) WithClientOrderId(id string) *BaseOrderOption {
	o.clientOrderId = id
	return o
}
func (o *BaseOrderOption) WithRelativeStopLoss(price int64) *BaseOrderOption {
	o.relativeStopLoss = price
	return o
}
func (o *BaseOrderOption) WithRelativeTakeProfit(price int64) *BaseOrderOption {
	o.relativeTakeProfit = price
	return o
}
func (o *BaseOrderOption) WithGuaranteedStopLoss(enabled bool) *BaseOrderOption {
	o.guaranteedStopLoss = enabled
	return o
}
func (o *BaseOrderOption) WithTrailingStopLoss(enabled bool) *BaseOrderOption {
	o.trailingStopLoss = enabled
	return o
}
func (o *BaseOrderOption) WithStopTriggerMethod(method openapi.ProtoOAOrderTriggerMethod) *BaseOrderOption {
	o.stopTriggerMethod = method
	return o
}

//////////////////
type AmendPositionSLTPOption struct {
	stopLoss              float64
	takeProfit            float64
	guaranteedStopLoss    bool
	trailingStopLoss      bool
	stopLossTriggerMethod openapi.ProtoOAOrderTriggerMethod
}

func (o *AmendPositionSLTPOption) WithStopLoss(price float64) *AmendPositionSLTPOption {
	o.stopLoss = price
	return o
}
func (o *AmendPositionSLTPOption) WithTakeProfit(price float64) *AmendPositionSLTPOption {
	o.takeProfit = price
	return o
}
func (o *AmendPositionSLTPOption) WithGuaranteedStopLoss(enabled bool) *AmendPositionSLTPOption {
	o.guaranteedStopLoss = enabled
	return o
}
func (o *AmendPositionSLTPOption) WithTrailingStopLoss(enabled bool) *AmendPositionSLTPOption {
	o.trailingStopLoss = enabled
	return o
}
func (o *AmendPositionSLTPOption) WithStopLossTriggerMethod(method openapi.ProtoOAOrderTriggerMethod) *AmendPositionSLTPOption {
	o.stopLossTriggerMethod = method
	return o
}
