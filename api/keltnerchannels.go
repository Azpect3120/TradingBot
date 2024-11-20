package api

import (
	"fmt"
	"math"
)

type KeltnerChannels struct {
	Source    Source  // Default: Close
	Length    int     // Default: 14
	Multipler float64 // Default: 2
	Upper     float64
	Basis     float64
	Lower     float64
}

// Creates a new instance of KeltnerChannels data structure.
// Values are set to default values, but can be changed by
// the caller.
func NewKeltnerChannels() *KeltnerChannels {
	return &KeltnerChannels{
		Source:    Close,
		Length:    14,
		Multipler: 2,
		Upper:     0.0,
		Basis:     0.0,
		Lower:     0.0,
	}
}

// Calculate the upper, basis, and lower keltner channels for the
// given bars. The values are stored in the KeltnerChannels
// structure. The calculation is based on the source, length, and
// multiplier values in the called structure. To calculate squeeze
// values, call this function multiple times with different
// multiplier values.
//
// WIP: The calculation does not change based on the source value.
func (kc *KeltnerChannels) Calculate(bars []Bar) {
	var (
		sum   float64 = 0.0
		avg   float64 = 0.0
		sumTR float64 = 0.0
		avgTR float64 = 0.0
	)

	// Calculate sum of the bars ATR
	var lastBar Bar = bars[len(bars)-(kc.Length)-1]
	for _, bar := range bars[len(bars)-(kc.Length) : len(bars)] {
		sum += bar.Close

		var TR float64 = math.Max(bar.High-bar.Low, math.Max(math.Abs(bar.High-lastBar.Close), math.Abs(bar.Low-lastBar.Close)))

		sumTR += TR
		lastBar = bar
	}

	// Set the average and basis value
	avgTR = sumTR / float64(kc.Length)
	avg = sum / float64(kc.Length)

	kc.Basis = avg

	// Calculate the upper and lower bands
	kc.Upper = avg + (avgTR * kc.Multipler)
	kc.Lower = avg - (avgTR * kc.Multipler)
}

// Return the string representation of the KeltnerChannels
// structure. The values are separated by a space. Starting with
// upper, basis, and lower.
func (kc *KeltnerChannels) String() string {
	return fmt.Sprintf("%f %f %f", kc.Upper, kc.Basis, kc.Lower)
}

// Return the string representation of the KeltnerChannels
// structure. The values are separated by a space. Starting with
// upper, basis, and lower. The values are formatted to the
// given precision.
func (kc *KeltnerChannels) StringFixed(precision int) string {
	format := fmt.Sprintf("%%.%df %%.%df %%.%df", precision, precision, precision)
	return fmt.Sprintf(format, kc.Upper, kc.Basis, kc.Lower)
}
