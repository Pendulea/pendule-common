package pcommon

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
	Dependencies               []AssetJSON   `json:"dependencies"`
	Arguments                  []string      `json:"arguments"`
}
