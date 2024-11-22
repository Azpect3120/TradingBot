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
	symPath := flag.String("f", "", "Path to the CSV file containing symbols.")
	symbol := flag.String("s", "", "Symbol to get data for. Do not use with -f.")
	lookback := flag.Int("l", 500, "Number of candles to look back.")
	distribution := flag.Bool("d", false, "Calculate the distribution of ratings. This must be used with -f to provided a list of symbols.")
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

		report := api.GenerateReport(*symbol, bars)
		fmt.Println(report.String())

		os.Exit(0)
	}

	// Distribution flag was provided
	if *distribution && *symPath != "" {
		rows, err := util.GetNamesFromCSV(*symPath)
		if err != nil {
			panic(err)
		}

		os.Exit(CalculateDistribution(rows, *lookback))
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

			report := api.GenerateReport(row, bars)
			if report.Rating.LongScore >= 50 || report.Rating.ShortScore >= 50 {
				fmt.Println(report.String())
			}
		}
		os.Exit(0)
	}
}

// Calculate the distribution of ratings for the symbols provided
// in the CSV file. The lookback is the number of candles to look back
// for the rating. The function will print out the distribution of
// ratings for the long and short direction. The function will return an
// int to be used as the exit code. This function will only fail if the
// API fails somehow.
func CalculateDistribution(rows []string, lookback int) int {
	// Store the results in an array, the index will store the ratings of
	// 0-10, 10-20, 20-30, 30-40, 40-50, 50-60, 60-70, 70-80, 80-90, 90-100
	var results_long [10]int = [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var results_short [10]int = [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	fmt.Printf("Generating distribution of ratings for %d symbols\n", len(rows))

	// Calculate the distribution of ratings
	for _, row := range rows {
		bars, err := util.GetBars(row, lookback)
		if err != nil {
			println(err.Error(), row)
			return 1
		}
		if len(bars) < 50 {
			continue
		}
		bars = bars[:len(bars)-1]

		api.GenerateReport(row, bars)
		// if rating.LongScore >= 0 {
		// 	results_long[int(rating.LongScore/10)]++
		// }
		// if rating.LongScore >= 0 {
		// 	results_short[int(rating.ShortScore/10)]++
		// }

		fmt.Printf("\rCalculating rating for $%-5s", row)
	}
	fmt.Printf(
		"\nLong Ratings Distribution\n[0-10]: %d\n[11-20]: %d\n[21-30]: %d\n[31-40]: %d\n[41-50]: %d\n[51-60]: %d\n[61-70]: %d\n[71-80]: %d\n[81-90]: %d\n[91-100]: %d\n",
		results_long[0], results_long[1], results_long[2], results_long[3], results_long[4], results_long[5], results_long[6], results_long[7], results_long[8], results_long[9],
	)
	fmt.Printf(
		"\nShort Ratings Distribution\n[0-10]: %d\n[11-20]: %d\n[21-30]: %d\n[31-40]: %d\n[41-50]: %d\n[51-60]: %d\n[61-70]: %d\n[71-80]: %d\n[81-90]: %d\n[91-100]: %d\n",
		results_short[0], results_short[1], results_short[2], results_short[3], results_short[4], results_short[5], results_short[6], results_short[7], results_short[8], results_short[9],
	)
	return 0
}
