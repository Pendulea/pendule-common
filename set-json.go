package pcommon

import (
	"encoding/json"
	"fmt"
)

type Consistency struct {
	Range     [2]TimeUnit `json:"range"`
	Timeframe int64       `json:"timeframe"`
}

type AssetJSON struct {
	ID                         AssetType     `json:"id"`
	Precision                  int8          `json:"precision"`
	Type                       DataType      `json:"type"`
	Consistencies              []Consistency `json:"consistencies"`
	ConsistencyMaxLookbackDays int           `json:"consistency_max_lookback_days"`
	SubAssets                  []AssetJSON   `json:"sub_assets"`
	Arguments                  []string      `json:"arguments"`
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
