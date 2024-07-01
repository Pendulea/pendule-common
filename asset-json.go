package pcommon

type Consistency struct {
	Range     [2]TimeUnit `json:"range"`
	Timeframe int64       `json:"timeframe"`
}

type AssetJSON struct {
	AddressString              AssetAddress       `json:"address_string"`
	Address                    AssetAddressParsed `json:"address"`
	ConsistencyMaxLookbackDays int                `json:"consistency_max_lookback_days"`
	Consistencies              []Consistency      `json:"consistencies"`
	DataType                   DataType           `json:"data_type"`
	Decimals                   int8               `json:"decimals"`
	MinDataDate                string             `json:"min_data_date"`
}
