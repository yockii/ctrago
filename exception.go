package ctrago

import "fmt"

var (
	ErrSymbolIdRequired      error = fmt.Errorf("symbolId is required")
	ErrAssetIdRequired       error = fmt.Errorf("assetId is required")
	ErrFromTimestampRequired error = fmt.Errorf("fromTimestamp is required")
	ErrToTimestampRequired   error = fmt.Errorf("toTimestamp is required")
	ErrTimestampRange        error = fmt.Errorf("timestamp range is invalid, it should be less than 7 days")
	ErrVolumeRequired        error = fmt.Errorf("volume is required")
	ErrPositionIdRequired    error = fmt.Errorf("positionId is required")
)
