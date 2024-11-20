package main

import (
	"fmt"
	"time"

	"github.com/Azpect3120/TradingBot/api"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

func main() {
	end := time.Now()
	start := end.Add(-500 * time.Hour)

	params := &chart.Params{
		Symbol:     "SU",
		Interval:   datetime.OneHour,
		IncludeExt: false,
		Start:      datetime.New(&start),
		End:        datetime.New(&end),
	}

	ch := chart.Get(params)

	var bars []api.Bar
	for ch.Next() {
		high, _ := ch.Bar().High.Float64()
		low, _ := ch.Bar().Low.Float64()
		open, _ := ch.Bar().Open.Float64()
		close_, _ := ch.Bar().Close.Float64()

		bar := api.Bar{
			High:      high,
			Low:       low,
			Open:      open,
			Close:     close_,
			Volume:    int64(ch.Bar().Volume),
			Timestamp: time.Unix(int64(ch.Bar().Timestamp), 0).UTC().Local(),
		}
		bars = append(bars, bar)
	}
	if err := ch.Err(); err != nil {
		fmt.Println(err)
	}

	// Remove the last bar because its fucked up
	bars = bars[:len(bars)-1]

	// for I, bar := range bars {
	// 	fmt.Printf("[%d] %s\n", I, bar.Timestamp.Format("2006-01-02 15:04:05"))
	// }

	sqz := api.NewSqueezePro()
	sqz.Calculate(bars)
	print(sqz.String())
}
