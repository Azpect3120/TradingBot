package util

import (
	"errors"
	"time"

	"github.com/Azpect3120/TradingBot/api"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

// Get the bars for a given symbol and number of hours. This
// function will return an array of Bar structs for the given
// symbol and number of hours. The interval will always be hourly
// since that is the use case and extended hours will be ignored.
// Hours must not be 0.
//
// This function will not clean the bars and remove the final
// "half bar." That should be done by the caller using basic
// array manipulation.
func GetBars(symbol string, hours int) ([]api.Bar, error) {
	// Cannot get 0 bars
	if hours == 0 {
		return nil, errors.New("Cannot get 0 bars")
	}

	// If the provided value is negative, set it back to positive
	// so it can be negated later.
	if hours < 0 {
		hours = -hours
	}

	// Calculate time range
	var end time.Time = time.Now()
	var start time.Time = end.Add(time.Duration(-hours) * time.Hour)

	params := &chart.Params{
		Symbol:     symbol,
		Interval:   datetime.OneHour,
		IncludeExt: false,
		Start:      datetime.New(&start),
		End:        datetime.New(&end),
	}

	var bars []api.Bar

	// Get the chart data
	ch := chart.Get(params)
	for ch.Next() {
		bar := parseBar(ch)
		bars = append(bars, bar)
	}
	if err := ch.Err(); err != nil {
		return nil, err
	}

	return bars, nil
}

// Convert the chart iterator to a Bar struct.
// This is just a convenience function for use in the
// GetBars function.
func parseBar(iter *chart.Iter) api.Bar {
	high, _ := iter.Bar().High.Float64()
	low, _ := iter.Bar().Low.Float64()
	open, _ := iter.Bar().Open.Float64()
	close_, _ := iter.Bar().Close.Float64()

	return api.Bar{
		High:      high,
		Low:       low,
		Open:      open,
		Close:     close_,
		Volume:    int64(iter.Bar().Volume),
		Timestamp: time.Unix(int64(iter.Bar().Timestamp), 0).UTC().Local(),
	}
}
