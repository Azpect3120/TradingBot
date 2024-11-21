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

	for _, row := range rows[:5] {
		bars, err := util.GetBars(row, 500)
		if err != nil {
			panic(err)
		}
		bars = bars[:len(bars)-1]

		print(row, ": ")
		api.CalculateRating(row, bars)
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
