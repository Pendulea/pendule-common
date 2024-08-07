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

	b := NewIndicatorDataBuilder(Asset.RSI, nil, []string{PERIOD_STRING}, -1)
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

func TestRSI2Builder(t *testing.T) {

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
	const COLUMN_NAME = "open"

	b := NewIndicatorDataBuilder(Asset.RSI2, nil, []string{COLUMN_NAME, PERIOD_STRING}, -1)
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

func TestSMABuilder(t *testing.T) {

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

	const PERIOD_STRING = "25"
	const PERIOD_INT = 25
	const COLUMN_NAME = "close"

	b := NewIndicatorDataBuilder(Asset.SMA, nil, []string{COLUMN_NAME, PERIOD_STRING}, 7)
	assert.Equal(t, []byte(nil), b.PrevState(), "PrevState should be nil")

	for i, unit := range units {
		point, err := b.ComputeUnsafe(unit)
		assert.Equal(t, nil, err, "Error should be nil")

		if i+1 < PERIOD_INT {
			assert.Equal(t, float64(-1), point.Value, "SMA value should be -1 before period is reached")
		} else {
			// Manually calculate the expected SMA value for validation
			sum := 0.0
			for j := i - PERIOD_INT + 1; j <= i; j++ {
				sum += btcPriceValues[j]
			}
			expectedSMA := Math.RoundFloat(sum/float64(PERIOD_INT), 7)
			assert.Equal(t, expectedSMA, point.Value, "SMA value should be equal to the expected SMA")
		}
	}
}

func TestEMABuilder(t *testing.T) {

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

	const PERIOD_STRING = "25"
	const PERIOD_INT = 25
	const COLUMN_NAME = "close"

	b := NewIndicatorDataBuilder(Asset.EMA, nil, []string{COLUMN_NAME, PERIOD_STRING}, 7)
	assert.Equal(t, []byte(nil), b.PrevState(), "PrevState should be nil")
	state := newEmptyMAState(PERIOD_INT)

	for i, unit := range units {
		point, err := b.ComputeUnsafe(unit)
		assert.Equal(t, nil, err, "Error should be nil")
		if i == 0 {
			state.EMA = unit.Close
			assert.Equal(t, unit.Close, point.Value, "First EMA value should be equal to the first unit close value")
		} else {
			// Manually calculate the expected EMA value for validation
			k := 2.0 / float64(PERIOD_INT+1)
			expectedEMA := unit.Close*k + state.EMA*(1-k)
			assert.Equal(t, Math.RoundFloat(expectedEMA, 7), point.Value, "EMA value should be equal to the expected EMA within a small epsilon")
			state.EMA = expectedEMA
		}
	}
}

func TestHMABuilder(t *testing.T) {

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

	const PERIOD_STRING = "9"
	const PERIOD_INT = 9
	const COLUMN_NAME = "close"

	b := NewIndicatorDataBuilder(Asset.HMA, nil, []string{COLUMN_NAME, PERIOD_STRING}, 7)
	assert.Equal(t, []byte(nil), b.PrevState(), "PrevState should be nil")

	for i, unit := range units {
		point, err := b.ComputeUnsafe(unit)
		assert.Equal(t, nil, err, "Error should be nil")
		if i <= PERIOD_INT {
			assert.Equal(t, float64(-1), point.Value, "HMA value should be -1 before period is reached")
		} else {
			assert.Greater(t, point.Value, 65000.00, "HMA value should be greater than 0")
		}
	}
}

func TestWMABuilder(t *testing.T) {

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

	const PERIOD_STRING = "20"
	const PERIOD_INT = 20
	const COLUMN_NAME = "close"

	b := NewIndicatorDataBuilder(Asset.WMA, nil, []string{COLUMN_NAME, PERIOD_STRING}, 7)
	assert.Equal(t, []byte(nil), b.PrevState(), "PrevState should be nil")

	for i, unit := range units {
		point, err := b.ComputeUnsafe(unit)
		assert.Equal(t, nil, err, "Error should be nil")
		if i+1 < PERIOD_INT {
			assert.Equal(t, float64(-1), point.Value, "WMA value should be -1 before period is reached")
		} else {
			assert.Greater(t, point.Value, 65000.00, "WMA value should be greater than 0")
		}
	}
}
