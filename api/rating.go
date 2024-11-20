package api

type Rating struct {
	Score   float64 // Score out of 100
	Squeeze Squeeze // Squeeze value
	Symbol  string
}

// Calculate the rating based on the bars passed into the
// function. The rating is based on many measures and can be
// found described in the function. The symbol passed into the
// function has no relevance, it is just for use in reporting
// later.
func CalculateRating(symbol string, bars []Bar) *Rating {
	rating := &Rating{
		Score:   0.0,
		Squeeze: SqueezeNone,
		Symbol:  symbol,
	}

	// Create SqueezePro instance and calculate squeeze
	sqz := NewSqueezePro(24)
	sqz.Calculate(bars)

	println(sqz.String())

	return rating
}

// Return the string representation of the rating.
//
// Implementation does not exist.
func (r *Rating) String() string {
	return ""
}
