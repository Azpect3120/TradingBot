package main

import (
	"github.com/Azpect3120/TradingBot/api"
	"github.com/Azpect3120/TradingBot/internal/util"
)

func main() {
	// end := time.Now()
	// start := end.Add(-500 * time.Hour)
	//
	// params := &chart.Params{
	// 	Symbol:     "SPY",
	// 	Interval:   datetime.OneHour,
	// 	IncludeExt: false,
	// 	Start:      datetime.New(&start),
	// 	End:        datetime.New(&end),
	// }
	//
	// ch := chart.Get(params)
	//
	// var bars []api.Bar
	// for ch.Next() {
	// 	high, _ := ch.Bar().High.Float64()
	// 	low, _ := ch.Bar().Low.Float64()
	// 	open, _ := ch.Bar().Open.Float64()
	// 	close_, _ := ch.Bar().Close.Float64()
	//
	// 	bar := api.Bar{
	// 		High:      high,
	// 		Low:       low,
	// 		Open:      open,
	// 		Close:     close_,
	// 		Volume:    int64(ch.Bar().Volume),
	// 		Timestamp: time.Unix(int64(ch.Bar().Timestamp), 0).UTC().Local(),
	// 	}
	// 	bars = append(bars, bar)
	// }
	// if err := ch.Err(); err != nil {
	// 	fmt.Println(err)
	// }

	bars, err := util.GetBars("SPY", 500)
	if err != nil {
		panic(err)
	}

	// Remove the last bar because its fucked up
	bars = bars[:len(bars)-1]

	for _, bar := range bars {
		println(bar.StringFixed(2))
	}

	sqz := api.NewSqueezePro()
	sqz.Calculate(bars)
	print(sqz.String())

	ma := api.NewMovingAverage(50)
	ma.Calculate(bars)
	println(ma.StringFixed(2))
}
