package api

import (
	"bytes"
	"strings"
	"text/template"
)

const TMPLPATH string = "./internal/templates/report.tmpl"

type Relation int

const (
	RelationNone  Relation = -1
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

	Squeeze         Squeeze   // Most recent squeeze value
	SqueezeHistory  []Squeeze // Historical squeeze values (past 14 bars)
	SqueezeLength   int       // Length of the squeeze
	SqueezeIncrease bool      // Squeeze is increasing
	SqueezeDecrease bool      // Squeeze is decreasing
	SqueezeConstant bool      // Squeeze is constant

	MA50Relation      Relation // Relation of the price to the 50 MA
	MA50CrossRelation Relation // 50 MA cross relation (most recent only) [set to -1 if no cross]
	MA50Increasing    bool     // 50 MA is increasing
	MA50Decreasing    bool     // 50 MA is decreasing
	MA50AboveHigh     bool     // 50 MA is above the high (past 21 bars) [for longs]
	MA50BelowLow      bool     // 50 MA is below the low (past 21 bars) [for shorts]
	MA50AboveLength   int      // 50 MA is above the high for this many bars
	MA50BelowLength   int      // 50 MA is below the low for this many bars

	MA9Relation       Relation // Relation of the price to the 9 MA
	MA9CrossRelation  Relation // 9 MA cross relation (most recent only) [set to -1 if no cross]
	MA9Increasing     bool     // 9 MA is increasing
	MA9Decreasing     bool     // 9 MA is decreasing
	MA9BetweenHighLow bool     // 9 MA is between the high and low (past 21 bars)
	MA9AboveLength    int      // 9 MA is above the high for this many bars
	MA9BelowLength    int      // 9 MA is below the low for this many bars
}

// Create a new report variable with the given symbol and rating.
// When the rating is changed, the report will be updated. SqueezeHistory
// is initialized as an empty slice, values should be appended to it.
// The first value will be the most recent squeeze.
func NewReport(symbol string, rating *Rating) *Report {
	return &Report{
		Symbol:            symbol,
		Rating:            rating,
		SqueezeHistory:    []Squeeze{},
		MA9CrossRelation:  RelationNone,
		MA50CrossRelation: RelationNone,
	}
}

// Print returns a string representation of the report.
// This function does not discriminate and will generate a
// result for any report provided, even shitty ones. The
// caller is expected to only call this function on desired
// reports.
func (r *Report) String() string {
	const tmp = `
  -----------------------------------
  | Report for ${{ .Symbol }}
  | ---------------------------------
  | Direction: {{ if eq .Direction 0 }}LONG{{ else }}SHORT{{ end }}
  | Rating: {{ if eq .Direction 0 }}{{ .Rating.LongScore }}{{ else }}{{ .Rating.ShortScore }}{{ end }}/100
  | ---------------------------------
  | Last Squeeze   : {{ if eq .Squeeze 2 }}WIDE{{ else if eq .Squeeze 3 }}NORMAL{{ else if eq .Squeeze 4 }}TIGHT{{ else if eq .Squeeze 5 }}VERY TIGHT{{ else }}NONE{{ end }}
  | Past Squeeze   : {{ range .SqueezeHistory }}{{ if eq . 2 }}W{{ else if eq . 3 }}N{{ else if eq . 4 }}T{{ else if eq . 5 }}V{{ else }}-{{ end }}{{ end }}
  | Squeeze Length : {{ .SqueezeLength }}
  | Squeeze Status : {{ if .SqueezeIncrease }}Tightening{{ else if .SqueezeDecrease }}Loosening{{ else if .SqueezeConstant }}Constant{{ else }}Unknown{{ end }}
  | 
  | Price to 50MA    : {{ if eq .MA50Relation 0 }}Above{{ else if eq .MA50Relation 1 }}Below{{ else }}None{{ end }}
  | Price cross 50MA : {{ if eq .MA50CrossRelation 0 }}Up{{ else if eq .MA50CrossRelation 1 }}Down{{ else }}--{{ end }}
  | 50MA Direction   : {{ if .MA50Increasing }}Increasing{{ else if .MA50Decreasing }}Decreasing{{ else }}None{{ end }}
  | {{ if eq .Direction 1 }}50MA Above High  : {{ if .MA50AboveHigh }}Yes{{ else }}No{{ end }}{{ else if eq .Direction 0 }}50MA Below Low   : {{ if .MA50BelowLow }}Yes{{ else }}No{{ end }}{{ end }}
  | {{ if eq .Direction 0 }}50MA Above Bars  : {{ .MA50AboveLength }}{{ else if eq .Direction 1}}50MA Below Bars  : {{ .MA50BelowLength }}{{ end }}
  | 
  | Price to 9MA     : {{ if eq .MA9Relation 0 }}Above{{ else if eq .MA9Relation 1 }}Below{{ else }}None{{ end }}
  | Price cross 9MA  : {{ if eq .MA9CrossRelation 0 }}Up{{ else if eq .MA9CrossRelation 1 }}Down{{ else }}--{{ end }}
  | 9MA Direction    : {{ if .MA9Increasing }}Increasing{{ else if .MA9Decreasing }}Decreasing{{ else }}None{{ end }}
  | 9MA Between H/L  : {{ if .MA9BetweenHighLow }}Yes{{ else }}No{{ end }}
  | {{ if eq .Direction 0 }}9MA Above Bars   : {{ .MA9AboveLength }}{{ else if eq .Direction 1 }}9MA Below Bars   : {{ .MA9BelowLength }}{{ end }}
  -----------------------------------
`
	tmpl, err := template.New("report").Parse(tmp)
	if err != nil {
		panic(err)
	}

	r.Symbol = strings.ToUpper(r.Symbol)

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, *r); err != nil {
		panic(err)
	}

	return buf.String()
}
