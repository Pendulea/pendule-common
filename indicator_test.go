package pcommon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIndicatorBuilder(t *testing.T) {

	btcPriceValues := []float64{
		65532.01, 65531.50, 65530.98, 65529.45, 65528.12,
		65527.65, 65526.78, 65525.89, 65524.47, 65523.55,
		65522.33, 66433.97, 65520.75, 65519.43, 65518.21,
		65517.99, 65516.88, 65515.65, 65514.42, 65513.18,
		65512.77, 65511.56, 65510.34, 65509.22, 65508.15,
		65507.03, 65506.89, 65505.77, 65504.66, 65503.54,
		65502.43, 65501.32, 65500.21, 65499.11, 65498.05,
		65497.92, 65496.81, 65495.71, 65494.62, 65493.53,
		65492.45, 65491.38, 65490.31, 65489.24, 65488.18,
		65487.12, 65486.07, 65485.02, 65484.96, 65483.91,
		65482.87, 65481.83, 65480.79, 65479.76, 65478.73,
		65477.70, 65476.68, 65475.66, 65474.64, 65473.63,
		65472.62, 65471.61, 65470.61, 65469.61, 65468.61,
		65467.62, 65466.63, 65465.64, 65464.66, 65463.68,
		65462.70, 65461.72, 65460.75, 65459.78, 65458.81,
		65457.85, 65456.89, 65455.93, 65454.98, 65454.02,
		65453.07, 65452.12, 65451.18, 65450.24, 65449.30,
		65448.37, 65447.44, 65446.51, 65445.58, 65444.66,
		65443.74, 65442.82, 65441.91, 65440.99, 65440.08,
		65439.18, 65438.27, 65437.37, 65436.47, 65435.58,
	}

	units := make(UnitTimeArray, 0)
	t0 := NewTimeUnit(1587607201)
	for i, price := range btcPriceValues {
		units = append(units, NewUnit(price).ToTime(t0.Add((time.Duration(i+1) * time.Second))))
	}

	const PERIOD_STRING = "14"
	const PERIOD_INT = 14

	b := NewIndicatorDataBuilder(Asset.RSI, nil, []string{PERIOD_STRING})
	assert.Equal(t, []byte(nil), b.PrevState(), "PrevState should be nil")
	p, err := b.ComputeUnsafe(units[0])
	assert.Equal(t, nil, err, "Error should be nil")
	state := b.cachedParsedState.(*rsiState)
	assert.Equal(t, units[0].Close, state.LastClose, "LastClose should be equal")
	assert.Equal(t, p.Value, -1.00, "Value should be equal")

	for i := 1; i < PERIOD_INT; i++ {
		_, err := b.ComputeUnsafe(units[i])
		assert.Equal(t, nil, err, "Error should be nil")
		assert.Equal(t, p.Value, -1.00, "Value should be equal")
		state := b.cachedParsedState.(*rsiState)
		assert.Equal(t, units[i].Close, state.LastClose, "LastClose should be equal")
	}
	p, err = b.ComputeUnsafe(units[PERIOD_INT])
	assert.Equal(t, nil, err, "Error should be nil")
	state = b.cachedParsedState.(*rsiState)
	assert.Equal(t, units[PERIOD_INT].Close, state.LastClose, "LastClose should be equal")
	assert.Equal(t, p.Value, 49.62440394539153, "Value should be equal")

}
