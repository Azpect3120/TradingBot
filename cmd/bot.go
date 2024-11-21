package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Azpect3120/TradingBot/api"
	"github.com/Azpect3120/TradingBot/internal/util"
)

func main() {
	// Create input flags for CLI
	symPath := flag.String("f", "", "Path to the CSV file containing symbols")
	symbol := flag.String("s", "", "Symbol to get data for. Do not use with -f")
	lookback := flag.Int("l", 500, "Number of days to look back")
	flag.Parse()

	// Symbol was provided
	if *symbol != "" {
		bars, err := util.GetBars(*symbol, *lookback)
		if err != nil {
			println(err.Error(), *symbol)
			os.Exit(1)
		}
		if len(bars) < 50 {
			println("Not enough bars were found", *symbol)
			os.Exit(1)
		}
		bars = bars[:len(bars)-1]

		rating := api.CalculateRating(*symbol, bars)
		fmt.Printf("Rating %s Long: %0.2f Short: %0.2f\n", *symbol, rating.LongScore, rating.ShortScore)

	}

	// Path was provided and no symbol
	if *symPath != "" && *symbol == "" {
		rows, err := util.GetNamesFromCSV(*symPath)
		if err != nil {
			panic(err)
		}

		for _, row := range rows {
			bars, err := util.GetBars(row, *lookback)
			if err != nil {
				println(err.Error(), row)
				continue
			}
			if len(bars) < 50 {
				println("Not enough bars were found", row)
				continue
			}
			bars = bars[:len(bars)-1]

			rating := api.CalculateRating(row, bars)
			// if rating.LongScore >= 50 || rating.ShortScore >= 50 {
			fmt.Printf("Rating %s Long: %0.2f Short: %0.2f\n", row, rating.LongScore, rating.ShortScore)
			// }
		}
	}

}
