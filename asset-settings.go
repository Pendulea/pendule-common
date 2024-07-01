package pcommon

type AssetSettings struct {
	Address     AssetAddressParsed `json:"address"`
	MinDataDate string             `json:"min_data_date"`
	Decimals    int8               `json:"decimals"`
}
