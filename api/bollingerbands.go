package api

import (
	"fmt"
	"math"
)

type BollingerBands struct {
	Source Source  // Default: Close
	Length int     // Default: 14
	StdDev float64 // Default: 2
	Upper  float64
	Basis  float64
	Lower  float64
}

// Creates a new instance of BollingerBands data structure.
// Values are set to default values, but can be changed by
// the caller.
func NewBollingerBands() *BollingerBands {
	return &BollingerBands{
		Source: Close,
		Length: 14,
		StdDev: 2,
		Upper:  0.0,
		Basis:  0.0,
		Lower:  0.0,
	}
}

// Calculate the upper, basis, and lower bollinger bands
// for the given bars. The values are stored in the BollingerBands
// structure. The calculation is based on the source, length, and
// standard deviation values in the called structure.
//
// WIP: The calculation does not change based on the source value.
func (bb *BollingerBands) Calculate(bars []Bar) {
	var (
		sum         float64 = 0.0
		avg         float64 = 0.0
		varianceSum float64 = 0.0
		variance    float64 = 0.0
		stddiv      float64 = 0.0
	)

	// Calculate sum of the bars
	for _, bar := range bars[len(bars)-(bb.Length) : len(bars)] {
		sum += bar.Close
	}

	// Set the average and basis value
	avg = sum / float64(bb.Length)
	bb.Basis = avg

	// Calculate variance sum of the bars
	for _, bar := range bars[len(bars)-(bb.Length) : len(bars)] {
		varianceSum += math.Pow(bar.Close-avg, 2)
	}

	// Calculate variance
	variance = varianceSum / float64(bb.Length)

	// Calculate standard deviation
	stddiv = math.Sqrt(variance)

	// Set upper and lower bands
	bb.Upper = avg + (bb.StdDev * stddiv)
	bb.Lower = avg - (bb.StdDev * stddiv)
}

// Return the string representation of the BollingerBands
// structure. The values are separated by a space. Starting with
// basis, upper, and lower.
func (bb *BollingerBands) String() string {
	return fmt.Sprintf("%f %f %f", bb.Basis, bb.Upper, bb.Lower)
}

// Return the string representation of the BollingerBands
// structure. The values are separated by a space. Starting with
// basis, upper, and lower. The values are formatted to the
// given precision.
func (bb *BollingerBands) StringFixed(precision int) string {
	format := fmt.Sprintf("%%.%df %%.%df %%.%df", precision, precision, precision)
	return fmt.Sprintf(format, bb.Basis, bb.Upper, bb.Lower)
}
