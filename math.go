package pcommon

import (
	"math"
	"sort"
)

type pmath struct{}

var Math pmath

// SafeMedian calculates the median of a slice of numbers
func (m pmath) SafeMedian(values []float64) float64 {
	length := len(values)
	if length == 0 {
		return 0
	}
	sort.Float64s(values)
	mid := length / 2
	if length%2 != 0 {
		return values[mid]
	}
	return (values[mid-1] + values[mid]) / 2
}

// SafeAverage calculates the average of a slice of numbers
func (m pmath) SafeAverage(values []float64) float64 {
	total := 0.0
	for _, value := range values {
		total += value
	}
	if len(values) == 0 {
		return 0
	}
	return total / float64(len(values))
}

func (m pmath) CalculateStandardDeviation(data []float64) float64 {
	if len(data) < 2 {
		return 0.0 // Standard deviation is not defined for one or zero elements.
	}

	// Step 1: Calculate the mean (average) of the data set
	mean := 0.0
	for _, value := range data {
		mean += value
	}
	mean /= float64(len(data))

	// Step 2: Calculate the variance (average of squared differences from the mean)
	variance := 0.0
	for _, value := range data {
		difference := value - mean
		squaredDifference := difference * difference
		variance += squaredDifference
	}
	variance /= float64(len(data) - 1) // Use N-1 for sample variance

	// Step 3: Standard deviation is the square root of the variance
	standardDeviation := math.Sqrt(variance)

	return standardDeviation
}

func (m pmath) RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
