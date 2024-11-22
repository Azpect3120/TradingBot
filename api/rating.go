package api

import (
	"math"
)

type Rating struct {
	score      float64 // Score out of 100
	LongScore  float64 // Long direction score
	ShortScore float64 // Short direction score
	Symbol     string
}

// Calculate the rating based on the bars passed into the
// function. The rating is based on many measures and can be
// found described in the function. The symbol passed into the
// function has no relevance, it is just for use in reporting
// later.
func GenerateReport(symbol string, bars []Bar) *Report {
	rating := &Rating{
		score:  0.0,
		Symbol: symbol,
	}

	report := NewReport(symbol, rating)

	// Create SqueezePro instance and calculate squeeze
	sqz := NewSqueezePro(len(bars))
	sqz.Calculate(bars)

	// Set the squeeze value in the report
	report.Squeeze = sqz.Squeeze[len(sqz.Squeeze)-1]

	// Generate score for recent squeeze
	switch sqz.Squeeze[len(sqz.Squeeze)-1] {
	case SqueezeVeryNarrow:
		rating.score += 15
	case SqueezeNarrow:
		rating.score += 10
	case SqueezeNormal:
		rating.score += 5
	case SqueezeWide:
		rating.score += 2.5
	case SqueezeNone, SqueezeUnknown:
		rating.score += 0
		return report // Failed case
	}

	// Calculate score for the historical squeeze

	// Loop over last 14 and count how many there are of each type.
	// Checks off a, b, c, d, and e
	var veryNarrow, narrow, normal, wide int = 0, 0, 0, 0
	for i := len(sqz.Squeeze) - 1; i >= len(sqz.Squeeze)-14; i-- {
		// Add squeeze to report
		report.SqueezeHistory = append(report.SqueezeHistory, sqz.Squeeze[i])

		// Calculate the score
		switch sqz.Squeeze[i] {
		case SqueezeVeryNarrow:
			if veryNarrow < 5 {
				veryNarrow++
				rating.score += 5
			}
		case SqueezeNarrow:
			if narrow < 5 {
				narrow++
				rating.score += 3
			}
		case SqueezeNormal:
			if normal < 5 {
				normal++
				rating.score += 1
			}
		case SqueezeWide:
			if wide < 5 {
				wide++
				rating.score += 0.5
			}
		case SqueezeNone, SqueezeUnknown:
			rating.score -= 1
		}
	}

	// Calculate the duration of the squeeze: checks off f
	var i int = len(sqz.Squeeze) - 1
	var count int = 0
	for sqz.Squeeze[i] > SqueezeNone {
		count++
		i--
	}
	var points int = int(math.Min(float64(count/7), 10.0))
	rating.score += float64(points) * 0.75

	// Add squeeze duration to the report
	report.SqueezeLength = count

	// Calculate if the squeeze is increasing: checks off g
	if sqz.Squeeze[len(sqz.Squeeze)-8] < sqz.Squeeze[len(sqz.Squeeze)-1] {
		rating.score += 5
		report.SqueezeIncrease = true
	}

	// Calculate if the squeeze is decreasing: checks off h
	if sqz.Squeeze[len(sqz.Squeeze)-8] > sqz.Squeeze[len(sqz.Squeeze)-1] {
		rating.score -= 5
		report.SqueezeDecrease = true
	}

	// Calculate if the squeeze is constant: checks off i
	if sqz.Squeeze[len(sqz.Squeeze)-8] == sqz.Squeeze[len(sqz.Squeeze)-1] {
		rating.score += 2.5
		report.SqueezeConstant = true
	}

	// Set both scores to the same value
	rating.LongScore = rating.score
	rating.ShortScore = rating.score

	// Calculate score for the 50HMA
	ma := NewMovingAverage(50)
	ma.Calculate(bars)

	// Check relation to the 50HMA on most recent bar: Checks off a and b
	if ma.Average[len(ma.Average)-1] > bars[len(bars)-1].Close {
		rating.ShortScore += 5
		rating.LongScore -= 2.5
		report.MA50Relation = RelationBelow
	}
	if ma.Average[len(ma.Average)-1] < bars[len(bars)-1].Close {
		rating.ShortScore -= 2.5
		rating.LongScore += 5
		report.MA50Relation = RelationAbove
	}

	// Check if price has crossed the 50HMA: Checks off c and d
	for i := len(bars) - 1; i >= len(bars)-7; i-- {
		// Cross was found. Now must determine which direction.
		// To do this, check if the crossing bar is higher than the 50HMA,
		// if so, the cross was up, and vise versa. Each crossing is scored.
		// This seems to make the most sense as a cross up and then down will
		// cancel out, which is the desired result for a neutral score.
		// This one only checks for crossing on green bars
		if bars[i].Close > ma.Average[i] && bars[i].Open < ma.Average[i] {
			if bars[i].Close > ma.Average[i] {
				rating.LongScore += 2.5
				rating.ShortScore -= 2.5

				if report.MA50CrossRelation == RelationNone {
					report.MA50CrossRelation = RelationAbove
				}
			}
			if bars[i].Close < ma.Average[i] {
				rating.LongScore -= 2.5
				rating.ShortScore += 2.5

				if report.MA50CrossRelation == RelationNone {
					report.MA50CrossRelation = RelationBelow
				}
			}
		}

		// This one only checks for crossing on red bars
		if bars[i].Close < ma.Average[i] && bars[i].Open > ma.Average[i] {
			if bars[i].Close > ma.Average[i] {
				rating.LongScore += 2.5
				rating.ShortScore -= 2.5

				if report.MA50CrossRelation == RelationNone {
					report.MA50CrossRelation = RelationAbove
				}
			}
			if bars[i].Close < ma.Average[i] {
				rating.LongScore -= 2.5
				rating.ShortScore += 2.5

				if report.MA50CrossRelation == RelationNone {
					report.MA50CrossRelation = RelationBelow
				}
			}
		}
	}

	// Calculate the duration above the 50HMA: checks off e (long)
	// Cannot run when the MA is 0, if the 50 is not calculated
	i = len(bars) - 1
	count = 0
	for bars[i].Close > ma.Average[i] && ma.Average[i] > 0 {
		count++
		i--
	}
	points = int(math.Min(float64(count/7), 5.0))
	rating.LongScore += float64(points) * 0.5

	// Add the duration above the 50HMA to the report
	report.MA50AboveLength = count

	// Calculate the duration below the 50HMA: checks off e (short)
	// Cannot run when the MA is 0, if the 50 is not calculated
	i = len(bars) - 1
	count = 0
	for bars[i].Close < ma.Average[i] && ma.Average[i] > 0 {
		count++
		i--
	}
	points = int(math.Min(float64(count/7), 5.0))
	rating.ShortScore += float64(points) * 0.5

	// Add the duration below the 50HMA to the report
	report.MA50BelowLength = count

	// Calculate 50HMA relation to 21H low/high: checks off f
	// Find low and high of the last 21 bars, based on low and high
	// not open and close. Which means it includes wicks.
	var low, high float64 = math.MaxFloat64, 0.0
	for i := len(bars) - 1; i >= len(bars)-21; i-- {
		low = math.Min(low, bars[i].Low)
		high = math.Max(high, bars[i].High)
	}

	// Checks off f (short)
	if ma.Average[len(ma.Average)-1] > high {
		rating.ShortScore += 2.5
		report.MA50AboveHigh = true
	}

	// Checks off f (long)
	if ma.Average[len(ma.Average)-1] < low {
		rating.LongScore += 2.5
		report.MA50BelowLow = true
	}

	// Calculate the 50HMA increasing or decreasing: checks off g
	i = 0
	var score int = 0
	for i < 5 {
		end := ma.Average[len(ma.Average)-(7*i)-1]
		start := ma.Average[len(ma.Average)-(7*(i+1))]

		if end > start {
			score++
		}
		if end < start {
			score--
		}

		i++
	}

	// Checks off g
	if score > 0 {
		rating.LongScore += float64(score) * 0.5
		report.MA50Increasing = true
	} else {
		rating.ShortScore += float64(-score) * 0.5
		report.MA50Decreasing = true
	}

	// Calculate score for the 9HMA
	ma.Length = 9
	ma.Calculate(bars)

	// Check relation to the 9HMA on most recent bar: Checks off a and b
	if ma.Average[len(ma.Average)-1] > bars[len(bars)-1].Close {
		rating.ShortScore += 3
		rating.LongScore -= 2
		report.MA9Relation = RelationBelow
	}
	if ma.Average[len(ma.Average)-1] < bars[len(bars)-1].Close {
		rating.ShortScore -= 2
		rating.LongScore += 3
		report.MA9Relation = RelationAbove
	}

	// Check if price has crossed the 9HMA: Checks off c and d
	for i := len(bars) - 1; i >= len(bars)-7; i-- {
		// Cross was found. Now must determine which direction.
		// To do this, check if the crossing bar is higher than the 50HMA,
		// if so, the cross was up, and vise versa. Each crossing is scored.
		// This seems to make the most sense as a cross up and then down will
		// cancel out, which is the desired result for a neutral score.
		// This one only checks for crossing on green bars
		if bars[i].Close > ma.Average[i] && bars[i].Open < ma.Average[i] {
			if bars[i].Close > ma.Average[i] {
				rating.LongScore += 2
				rating.ShortScore -= 2

				if report.MA9CrossRelation == RelationNone {
					report.MA9CrossRelation = RelationAbove
				}
			}
			if bars[i].Close < ma.Average[i] {
				rating.LongScore -= 2
				rating.ShortScore += 2

				if report.MA9CrossRelation == RelationNone {
					report.MA9CrossRelation = RelationBelow
				}
			}
		}

		// This one only checks for crossing on red bars
		if bars[i].Close < ma.Average[i] && bars[i].Open > ma.Average[i] {
			if bars[i].Close > ma.Average[i] {
				rating.LongScore += 2
				rating.ShortScore -= 2

				if report.MA9CrossRelation == RelationNone {
					report.MA9CrossRelation = RelationAbove
				}
			}
			if bars[i].Close < ma.Average[i] {
				rating.LongScore -= 2
				rating.ShortScore += 2

				if report.MA9CrossRelation == RelationNone {
					report.MA9CrossRelation = RelationBelow
				}
			}
		}
	}

	// Calculate is the 9HMA is between the highs and lows of the past 7 candles
	// Checks off e (long) and e (short)
	low, high = math.MaxFloat64, 0.0
	for i := len(bars) - 1; i >= len(bars)-7; i-- {
		low = math.Min(low, bars[i].Low)
		high = math.Max(high, bars[i].High)
	}
	if ma.Average[len(ma.Average)-1] > low && ma.Average[len(ma.Average)-1] < high {
		rating.LongScore += 2
		rating.ShortScore += 2
		report.MA9BetweenHighLow = true
	}

	// Calculate the duration above the 9HMA: checks off f (long)
	// Cannot run when the MA is 0, if the 50 is not calculated
	i = len(bars) - 1
	count = 0
	for bars[i].Close > ma.Average[i] && ma.Average[i] > 0 {
		count++
		i--
	}
	points = int(math.Min(float64(count/7), 3.0))
	rating.LongScore += float64(points) * 0.5

	// Add the duration above the 9HMA to the report
	report.MA9AboveLength = count

	// Calculate the duration below the 9HMA: checks off f (short)
	// Cannot run when the MA is 0, if the 50 is not calculated
	i = len(bars) - 1
	count = 0
	for bars[i].Close < ma.Average[i] && ma.Average[i] > 0 {
		count++
		i--
	}
	points = int(math.Min(float64(count/7), 3.0))
	rating.ShortScore += float64(points) * 0.5

	// Add the duration below the 9HMA to the report
	report.MA9BelowLength = count

	// Calculate the 50HMA increasing or decreasing: checks off g
	i = 0
	score = 0
	for i < 3 {
		end := ma.Average[len(ma.Average)-(7*i)-1]
		start := ma.Average[len(ma.Average)-(7*(i+1))]

		if end > start {
			score++
		}
		if end < start {
			score--
		}

		i++
	}

	// Checks off g
	if score > 0 {
		rating.LongScore += float64(score) * 0.5
		report.MA9Increasing = true
	} else {
		rating.ShortScore += float64(-score) * 0.5
		report.MA9Decreasing = true
	}

	// Determine the direction of the report
	if rating.LongScore > rating.ShortScore {
		report.Direction = DirectionLong
	} else {
		report.Direction = DirectionShort
	}

	return report
}

// Return the string representation of the rating.
//
// Implementation does not exist.
func (r *Rating) String() string {
	return ""
}
