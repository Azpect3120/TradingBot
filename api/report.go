package api

type Relation int

const (
	RelationAbove Relation = 0
	RelationBelow Relation = 1
)

type Direction int

const (
	DirectionLong  Direction = 0
	DirectionShort Direction = 1
)

// Report structure stores the data required to generate a report.
type Report struct {
	Symbol string  // Symbol of the report
	Rating *Rating // Rating generated for the report

	Direction Direction // Direction of the report

	Squeeze         Squeeze // Most recent squeeze value
	SqueezeLength   int     // Length of the squeeze
	SqueezeIncrease bool    // Squeeze is increasing
	SqueezeConstant bool    // Squeeze is constant

	MA50Relation      Relation // Relation of the price to the 50 MA
	MA50CrossRelation Relation // 50 MA cross relation (most recent only)
	MA50Increasing    bool     // 50 MA is increasing
	MA50Decreasing    bool     // 50 MA is decreasing
	MA50AboveHigh     bool     // 50 MA is above the high (last 21 bars) [for longs]
	MA50BelowLow      bool     // 50 MA is below the low (last 21 bars) [for shorts]

	MA9Relation       Relation // Relation of the price to the 9 MA
	MA9CrossRelation  Relation // 9 MA cross relation (most recent only)
	MA9Increasing     bool     // 9 MA is increasing
	MA9Decreasing     bool     // 9 MA is decreasing
	MA9BetweenHighLow bool     // 9 MA is between the high and low (last 7 bars)
}

// Create a new report variable with the given symbol and rating.
// When the rating is changed, the report will be updated.
func NewReport(symbol string, rating *Rating) *Report {
	return &Report{
		Symbol: symbol,
		Rating: rating,
	}
}
