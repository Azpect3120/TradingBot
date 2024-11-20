package api

import (
	"fmt"
	"strings"
)

type MovingAverage struct {
	Souce  Source // Default: close
	Length int
	Values []float64 // 0 where no value can be calculated
}

// Creates a new instance of MovingAverage data structure.
// Values are set to default values, but can be changed by
// the caller. The values are stored in a slice of float64
// values. The length of the slice is not set so it should
// not be accessed until values are added to the slice
// with the Calculate function.
func NewMovingAverage(length int) *MovingAverage {
	return &MovingAverage{
		Souce:  Close,
		Length: length,
		Values: []float64{},
	}
}

// Calculate the values for the moving average based on the
// bars passed in. The values are stored in the Values slice
// of the MovingAverage structure. The first value in the slice
// is the least recent value, and the last value is the most
// recent value. That way they are in the same order as the
// bars passed in and be accessed in the same order.
//
// WIP: The calculation does not change based on the source value.
func (ma *MovingAverage) Calculate(bars []Bar) {
	// Reset values
	ma.Values = []float64{}

	// Storing working sum
	var sum float64 = 0.0

	// Generate the sum of the first bars in the slice.
	// Then a window can be used. Fill the values with 0
	// since no value can be calculated yet.
	for i := 0; i < ma.Length; i++ {
		sum += bars[i].Close
		ma.Values = append(ma.Values, 0.0)
	}

	// Fill the values starting at the length of the slice.
	// A window is used to calculate the moving average
	// linearly.
	for i := ma.Length; i < len(bars); i++ {
		// Calculate the working sum
		sum += bars[i].Close
		sum -= bars[i-ma.Length].Close

		ma.Values = append(ma.Values, (sum / float64(ma.Length)))
	}
}

// String returns a string representation of the MovingAverage
// structure. The string representation is a comma separated
// list of values, starting with the least recent value.
func (ma *MovingAverage) String() string {
	var values []string

	for _, value := range ma.Values {
		values = append(values, fmt.Sprintf("%f", value))
	}

	return strings.Join(values, ", ")
}

// String returns a string representation of the MovingAverage
// structure. The string representation is a comma separated
// list of values, starting with the least recent value. The values
// are rounded to the given precision.
func (ma *MovingAverage) StringFixed(precision int) string {
	var values []string

	var format string = fmt.Sprintf("%%.%df", precision)
	for _, value := range ma.Values {
		values = append(values, fmt.Sprintf(format, value))
	}

	return strings.Join(values, ", ")
}
