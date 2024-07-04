package pcommon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountDivisionsTo(t *testing.T) {

	const BTC_PRICE = 60500.00
	const ETH_PRICE = 4000.00
	const LTC_PRICE = 80.00
	const CTSI_PRICE = 0.15
	const INJ_PRICE = 20.00
	const USDT_PRICE = 1.00
	const MATIC_PRICE = 0.60
	const OCEAN_PRICE = 0.50

	assert.Equal(t, int8(7), countDivisionsTo(BTC_PRICE, 0.01))
	assert.Equal(t, int8(6), countDivisionsTo(ETH_PRICE, 0.01))
	assert.Equal(t, int8(4), countDivisionsTo(LTC_PRICE, 0.01))
	assert.Equal(t, int8(2), countDivisionsTo(CTSI_PRICE, 0.01))
	assert.Equal(t, int8(4), countDivisionsTo(INJ_PRICE, 0.01))
	assert.Equal(t, int8(3), countDivisionsTo(USDT_PRICE, 0.01))
	assert.Equal(t, int8(2), countDivisionsTo(MATIC_PRICE, 0.01))
	assert.Equal(t, int8(2), countDivisionsTo(OCEAN_PRICE, 0.01))
}

func TestPriceDecimals(t *testing.T) {

	const BTC_PRICE = 60500.00
	const ETH_PRICE = 4000.00
	const LTC_PRICE = 80.00
	const CTSI_PRICE = 0.15
	const INJ_PRICE = 20.00
	const USDT_PRICE = 1.00
	const MATIC_PRICE = 0.60
	const OCEAN_PRICE = 0.50

	assert.Equal(t, int8(2), priceDecimals(BTC_PRICE))
	assert.Equal(t, int8(2), priceDecimals(ETH_PRICE))
	assert.Equal(t, int8(3), priceDecimals(LTC_PRICE))
	assert.Equal(t, int8(5), priceDecimals(CTSI_PRICE))
	assert.Equal(t, int8(3), priceDecimals(INJ_PRICE))
	assert.Equal(t, int8(4), priceDecimals(USDT_PRICE))
	assert.Equal(t, int8(5), priceDecimals(MATIC_PRICE))
	assert.Equal(t, int8(5), priceDecimals(OCEAN_PRICE))
}
