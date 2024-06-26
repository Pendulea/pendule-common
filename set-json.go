package pcommon

import (
	"encoding/json"
	"fmt"
)

type AssetJSON struct {
	ID                         AssetType   `json:"id"`
	Precision                  int8        `json:"precision"`
	Type                       DataType    `json:"type"`
	ConsistencyRange           [2]TimeUnit `json:"consistency_range"`
	ConsistencyMaxLookbackDays int         `json:"consistency_max_lookback_days"`
	Timeframe                  int64       `json:"timeframe"` // in milliseconds
	SubAssets                  []AssetJSON `json:"sub_assets"`
}

type SetJSON struct {
	Settings SetSettings `json:"settings"`
	Size     int64       `json:"size"`
	Assets   []AssetJSON `json:"assets"`
}

func (s SetJSON) PrintJSON() {
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
