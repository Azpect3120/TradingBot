package main

import (
	"github.com/Azpect3120/TradingBot/api"
	"github.com/Azpect3120/TradingBot/internal/util"
)

func main() {
	rows, err := util.GetNamesFromCSV("./data/large_mega_cap.csv")
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		bars, err := util.GetBars(row, 500)
		if err != nil {
			// Failed on these names: need to find a fix
			// BRK/A: changed to BRK-A
			// BRK/B: changed to BRK-B
			println(err.Error(), row)
			continue
		}
		if len(bars) == 0 {
			println("No bars found for ", row)
			continue
		}
		bars = bars[:len(bars)-1]

		api.CalculateRating(row, bars)
		println("Rating ", row)
	}

}

func test() {
	bars, err := util.GetBars("SPY", 500)
	if err != nil {
		panic(err)
	}

	// Remove the last bar because its fucked up
	// When running during market hours without this,
	// the last bar is incomplete and causes the program
	// to panic.
	bars = bars[:len(bars)-1]

	for _, bar := range bars {
		println(bar.StringFixed(2))
	}

	sqz := api.NewSqueezePro(15)
	sqz.Calculate(bars)
	print(sqz.String())

	ma := api.NewMovingAverage(50)
	ma.Calculate(bars)
	println(ma.StringFixed(2))

}
