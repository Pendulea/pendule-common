package pcommon

import "time"

type Consistency struct {
	Range     [2]TimeUnit `json:"range"`
	Timeframe int64       `json:"timeframe"`
	MinValue  float64     `json:"min_value"`
	MaxValue  float64     `json:"max_value"`
}

type AssetJSON struct {
	AddressString              AssetAddress           `json:"address_string"`
	Address                    AssetAddressParsedJSON `json:"address"`
	ConsistencyMaxLookbackDays int                    `json:"consistency_max_lookback_days"`
	Consistencies              []Consistency          `json:"consistencies"`
	DataType                   DataType               `json:"data_type"`
	Decimals                   int8                   `json:"decimals"`
	MinDataDate                string                 `json:"min_data_date"`
	LastReadTime               TimeUnit               `json:"last_read_time"`
}

func (a AssetJSON) FindConsistencyByTimeframe(timeframe time.Duration) *Consistency {
	for _, c := range a.Consistencies {
		if timeframe == time.Duration(c.Timeframe)*time.Millisecond {
			return &c
		}
	}
	return nil
}
