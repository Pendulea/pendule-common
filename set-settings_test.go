package pcommon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetJSON(t *testing.T) {
	const jsonDATA0 = `{
		"id": ["CTSI", "USDT"],
		"assets": [{
            "address": {
                "asset_type": "spot_price",
                "dependencies": [],
                "arguments": []
            },
			"min_data_date": "2020-05-05",
			"decimals": 3
		},
		{
            "address": {
                "asset_type": "spot_volume",
                "dependencies": [],
                "arguments": []
            },
			"min_data_date": "2020-05-05",
			"decimals": 3
		}],
		"settings": {
			"binance": 1
        }
	}`
	set := SetSettings{}
	err := json.Unmarshal([]byte(jsonDATA0), &set)
	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, set.IsValid(), nil, "Set should be valid")
	assert.Equal(t, set.IsBinancePair(), nil, "Set should be a binance pair")

	set.Assets = append(set.Assets, AssetSettings{
		Address: AssetAddressParsedWithoutSetID{
			AssetType:    Asset.RSI,
			Dependencies: []AssetAddress{set.Assets[0].Address.AddSetID(set.ID).BuildAddress()},
			Arguments:    []string{"14"},
		},
		MinDataDate: "2020-05-05",
	})
	assert.Equal(t, set.IsValid(), nil, "Set should be valid")
	assert.Equal(t, set.IsBinancePair(), nil, "Set should be a binance pair")
	set.Assets = append(set.Assets, AssetSettings{
		Address: AssetAddressParsedWithoutSetID{
			AssetType:    Asset.RSI,
			Dependencies: []AssetAddress{set.Assets[0].Address.AddSetID(set.ID).BuildAddress()},
			Arguments:    []string{"14", "14"},
		},
	})
	assert.NotEqual(t, set.IsValid(), nil, "Set should be invalid")
	set.Assets = set.Assets[:3]

	set.Assets = append(set.Assets, AssetSettings{
		Address: AssetAddressParsedWithoutSetID{
			AssetType:    Asset.RSI2,
			Dependencies: []AssetAddress{set.Assets[0].Address.AddSetID(set.ID).BuildAddress()},
			Arguments:    []string{"close", "14"},
		},
		MinDataDate: "2020-05-05",
	})
	assert.Equal(t, set.IsValid(), nil, "Set should be valid")
	assert.Equal(t, set.IsBinancePair(), nil, "Set should be a binance pair")

}
