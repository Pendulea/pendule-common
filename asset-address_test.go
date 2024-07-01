package pcommon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildAssetAddress(t *testing.T) {
	a := AssetAddressParsed{
		SetID:        []string{"ctsi", "usdt"},
		AssetType:    Asset.SPOT_PRICE,
		Dependencies: nil,
		Arguments:    nil,
	}
	aAddr := a.BuildAddress()
	assert.Equal(t, string(aAddr), "ctsi_usdt;spot_price;[];")
	aParsed, err := aAddr.Parse()
	assert.Nil(t, err)
	assert.Equal(t, a, *aParsed)

	a.BuildAddress()
	b := AssetAddressParsed{
		SetID:        []string{"ctsi", "usdt"},
		AssetType:    Asset.SPOT_VOLUME,
		Dependencies: nil,
		Arguments:    nil,
	}
	bAddr := b.BuildAddress()
	assert.Equal(t, string(bAddr), "ctsi_usdt;spot_volume;[];")
	bParsed, err := bAddr.Parse()
	assert.Nil(t, err)
	assert.Equal(t, b, *bParsed)

	c := AssetAddressParsed{
		SetID:     []string{"ctsi", "usdt"},
		AssetType: AssetType("vwap"),
		Dependencies: []AssetAddress{
			a.BuildAddress(),
			b.BuildAddress(),
		},
		Arguments: []string{"14"},
	}
	cAddr := c.BuildAddress()
	assert.Equal(t, string(cAddr), "ctsi_usdt;vwap;[ctsi_usdt;spot_price;[];=ctsi_usdt;spot_volume;[];];14")

	cParsed, err := cAddr.Parse()
	assert.Nil(t, err)
	assert.Equal(t, c, *cParsed)

	rsi := AssetAddressParsed{
		SetID:        []string{"ctsi", "usdt"},
		AssetType:    Asset.RSI,
		Dependencies: []AssetAddress{c.BuildAddress()},
		Arguments:    []string{"21"},
	}

	rsiAddr := rsi.BuildAddress()
	assert.Equal(t, string(rsiAddr), "ctsi_usdt;rsi;[ctsi_usdt;vwap;[ctsi_usdt;spot_price;[];=ctsi_usdt;spot_volume;[];];14];21")
	rsiParsed, err := rsiAddr.Parse()
	assert.Nil(t, err)
	assert.Equal(t, rsi, *rsiParsed)

	ema := AssetAddressParsed{
		SetID:        []string{"ctsi", "usdt"},
		AssetType:    AssetType("ema"),
		Dependencies: []AssetAddress{rsi.BuildAddress()},
		Arguments:    []string{"50"},
	}

	emaAddr := ema.BuildAddress()
	assert.Equal(t, string(emaAddr), "ctsi_usdt;ema;[ctsi_usdt;rsi;[ctsi_usdt;vwap;[ctsi_usdt;spot_price;[];=ctsi_usdt;spot_volume;[];];14];21];50")
	emaParsed, err := emaAddr.Parse()
	assert.Nil(t, err)
	assert.Equal(t, ema, *emaParsed)

	blabla := AssetAddressParsed{
		SetID:     []string{"general", "cryto"},
		AssetType: AssetType("blabla"),
		Dependencies: []AssetAddress{
			rsi.BuildAddress(),
			ema.BuildAddress(),
		},
		Arguments: []string{"12", "30"},
	}

	blablaAddr := blabla.BuildAddress()
	assert.Equal(t, string(blablaAddr), "general_cryto;blabla;[ctsi_usdt;rsi;[ctsi_usdt;vwap;[ctsi_usdt;spot_price;[];=ctsi_usdt;spot_volume;[];];14];21=ctsi_usdt;ema;[ctsi_usdt;rsi;[ctsi_usdt;vwap;[ctsi_usdt;spot_price;[];=ctsi_usdt;spot_volume;[];];14];21];50];12_30")
	blablaParsed, err := blablaAddr.Parse()
	assert.Nil(t, err)
	assert.Equal(t, blabla, *blablaParsed)
}
