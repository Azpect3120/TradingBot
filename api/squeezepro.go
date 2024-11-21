package api

import (
	"slices"
	"strings"
)

type SqueezePro struct {
	Length          int // Default: 14
	Thresholds      Thresholds
	BollingerBands  *BollingerBands
	KeltnerChannels *KeltnerChannels
	Squeeze         []Squeeze
	SqueezeCount    int
}

type Thresholds struct {
	Wide       float64 // Default: 1.8
	Normal     float64 // Default: 1.25
	Narrow     float64 // Default: 0.9
	VeryNarrow float64 // Default: 0.75
}

type Squeeze string

const (
	SqueezeWide       Squeeze = "Wide"
	SqueezeNormal     Squeeze = "Normal"
	SqueezeNarrow     Squeeze = "Narrow"
	SqueezeVeryNarrow Squeeze = "VeryNarrow"
	SqueezeNone       Squeeze = "None"
	SqueezeUnknown    Squeeze = "Unknown"
)

// Creates a new instance of SqueezePro data structure.
// Values are set to default values, but can be changed by
// the caller. The BollingerBands and KeltnerChannels
// structures are also created with default values.
func NewSqueezePro(size int) *SqueezePro {
	return &SqueezePro{
		Length: 14,
		Thresholds: Thresholds{
			Wide:       1.8,
			Normal:     1.25,
			Narrow:     0.9,
			VeryNarrow: 0.75,
		},
		BollingerBands:  NewBollingerBands(),
		KeltnerChannels: NewKeltnerChannels(),
		SqueezeCount:    size,
		Squeeze:         make([]Squeeze, size),
	}
}

// Calculate the squeeze pro values for the given bars.
// The values are stored in the SqueezePro structure.
// The calculation is based on the length and thresholds
// values in the called structure. The BollingerBands and
// KeltnerChannels structures are also calculated for use
// in the squeeze pro calculation. These are stored in
// the same order as the bars! They can be accessed using
// the same index as the provided bar and the value will
// represent the squeeze at that bar.
func (sqz *SqueezePro) Calculate(bars []Bar) {
	// N is based on the SqueezeCount value
	for n := 0; n < sqz.SqueezeCount; n++ {
		if n != 0 {
			bars = bars[:len(bars)-1]
		}

		// Ensure the bars are large enough to calculate the squeeze,
		// otherwise skip the calculation.
		if (len(bars) - sqz.KeltnerChannels.Length) <= 0 {
			sqz.Squeeze[n] = SqueezeUnknown
			continue
		}

		// Calculate the BB values
		sqz.BollingerBands.Calculate(bars)

		// Calculate very narrow squeeze
		sqz.KeltnerChannels.Multipler = sqz.Thresholds.VeryNarrow
		sqz.KeltnerChannels.Calculate(bars)
		if (sqz.BollingerBands.Lower >= sqz.KeltnerChannels.Lower) && (sqz.BollingerBands.Upper <= sqz.KeltnerChannels.Upper) {
			sqz.Squeeze[n] = SqueezeVeryNarrow
			continue
		}

		// Calculate narrow squeeze
		sqz.KeltnerChannels.Multipler = sqz.Thresholds.Narrow
		sqz.KeltnerChannels.Calculate(bars)
		if (sqz.BollingerBands.Lower >= sqz.KeltnerChannels.Lower) && (sqz.BollingerBands.Upper <= sqz.KeltnerChannels.Upper) {
			sqz.Squeeze[n] = SqueezeNarrow
			continue
		}

		// Calculate normal squeeze
		sqz.KeltnerChannels.Multipler = sqz.Thresholds.Normal
		sqz.KeltnerChannels.Calculate(bars)
		if (sqz.BollingerBands.Lower >= sqz.KeltnerChannels.Lower) && (sqz.BollingerBands.Upper <= sqz.KeltnerChannels.Upper) {
			sqz.Squeeze[n] = SqueezeNormal
			continue
		}

		// Calculate wide squeeze
		sqz.KeltnerChannels.Multipler = sqz.Thresholds.Wide
		sqz.KeltnerChannels.Calculate(bars)
		if (sqz.BollingerBands.Lower >= sqz.KeltnerChannels.Lower) && (sqz.BollingerBands.Upper <= sqz.KeltnerChannels.Upper) {
			sqz.Squeeze[n] = SqueezeWide
			continue
		}

		// No Squeeze
		sqz.Squeeze[n] = SqueezeNone
	}

	// Reverse the Squeeze array
	slices.Reverse(sqz.Squeeze)
}

// Return the string representation of the SqueezePro data.
// The string is the Squeeze value from the enum. Only displays
// the first value in the array. The first value is the most
// recent squeeze value.
func (sqz *SqueezePro) String() string {
	var result []string
	for _, s := range sqz.Squeeze {
		switch s {
		case SqueezeVeryNarrow:
			result = append(result, "Very Narrow (blue)")
		case SqueezeNarrow:
			result = append(result, "Narrow (yellow)")
		case SqueezeNormal:
			result = append(result, "Normal (red)")
		case SqueezeWide:
			result = append(result, "Wide (orange)")
		case SqueezeNone:
			result = append(result, "None (green)")
		case SqueezeUnknown:
			result = append(result, "Unknown (black)")
		}
	}

	return strings.Join(result, ",  ")
}

// Set the squeeze count for the SqueezePro structure.
// This is required to set the size of the Squeeze array.
func (sqz *SqueezePro) SetSqueezeCount(n int) {
	sqz.SqueezeCount = n
	sqz.Squeeze = make([]Squeeze, n)
}
