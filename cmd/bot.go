package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/Azpect3120/TradingBot/api"
	"github.com/Azpect3120/TradingBot/internal/util"
)

func main() {
	// Create input flags for CLI
	symbol := flag.String("s", "", "Symbol to get data for. Do not use with -f.")
	lookback := flag.Int("l", 500, "Number of candles to look back.")
	history := flag.Int("h", 0, "Number of days back to look for data.")
	distribution := flag.Bool("d", false, "Calculate the distribution of ratings. This must be used with -f to provided a list of symbols.")
	testing := flag.Bool("t", false, "Run the tests for the API.")
	flag.Parse()

	if *testing {
		tickers, err := api.GetTickers()
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nTicker count: %d\n", len(tickers))

		var count int
		for _, ticker := range tickers {
			vol, err := strconv.ParseInt(ticker.Volume, 10, 64)
			if err != nil {
				panic(err)
			}
			if vol >= 1000000 {
				fmt.Println(ticker.Symbol, ": ", ticker.LastSale)
				count++
			}
		}

		fmt.Printf("\nTicker count greater than 1M volume: %d\n", count)
	}

	// Symbol was provided
	if *symbol != "" {
		// bars, err := util.GetBars(*symbol, *lookback, *history)
		// if err != nil {
		// 	println(err.Error(), *symbol)
		// 	os.Exit(1)
		// }
		// if len(bars) < 50 {
		// 	println("Not enough bars were found", *symbol)
		// 	os.Exit(1)
		// }
		// bars = bars[:len(bars)-1]
		//
		// report := api.GenerateReport(*symbol, bars)
		// fmt.Println(report.String())

		os.Exit(0)
	}

	// Distribution flag was provided
	if *distribution {
		os.Exit(CalculateDistribution(*lookback, *history))
	}

	// Path was provided and no symbol
	if *symbol == "" {
		tickers, err := api.GetTickers()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, ticker := range tickers {
			marketCap, _ := strconv.ParseFloat(ticker.MarketCap, 64)
			volume, _ := strconv.ParseInt(ticker.Volume, 10, 64)

			// Skipping symbols with volume less than 1M or market cap less than 10B
			if volume < 1000000 || marketCap < 10000000000.00 {
				continue
			}

			bars, err := util.GetBars(ticker.Symbol, *lookback, *history)
			if err != nil {
				fmt.Printf("Failed to get bars for %s. Error: %s\n", ticker.Symbol, err.Error())
				continue
			}
			if len(bars) < 50 {
				println("Not enough bars were found", ticker.Symbol)
				continue
			}
			bars = bars[:len(bars)-1]

			report := api.GenerateReport(&ticker, bars)
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
// API fails somehow. All the names on the market will be queried. The
// names are only stored based on which value is bigger. A short OR long
// is stored for each name. Unless the ratings are the same, in which both
// are stored.
func CalculateDistribution(lookback, history int) int {
	// Store the results in an array, the index will store the ratings of
	// 0-10, 10-20, 20-30, 30-40, 40-50, 50-60, 60-70, 70-80, 80-90, 90-100
	var results_long [10]int = [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var results_short [10]int = [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// Store the failed symbols
	var failed []string

	tickers, err := api.GetTickers()
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	fmt.Printf("Generating oistribution of ratings for %d symbols\n", len(tickers))

	// Calculate the distribution of ratings
	for _, ticker := range tickers {
		bars, err := util.GetBars(ticker.Symbol, lookback, history)
		if err != nil {
			failed = append(failed, ticker.Symbol)
			continue
		}
		if len(bars) < 50 {
			continue
		}
		bars = bars[:len(bars)-1]

		report := api.GenerateReport(&ticker, bars)

		if report.Rating.LongScore >= 0 && report.Rating.LongScore >= report.Rating.ShortScore {
			results_long[int(report.Rating.LongScore/10)]++
		}
		if report.Rating.LongScore >= 0 && report.Rating.LongScore <= report.Rating.ShortScore {
			results_short[int(report.Rating.ShortScore/10)]++
		}

		fmt.Printf("\rCalculating rating for $%-5s", ticker.Symbol)
	}
	fmt.Printf(
		"\nLong Ratings Distribution\n[0-10]: %d\n[11-20]: %d\n[21-30]: %d\n[31-40]: %d\n[41-50]: %d\n[51-60]: %d\n[61-70]: %d\n[71-80]: %d\n[81-90]: %d\n[91-100]: %d\n",
		results_long[0], results_long[1], results_long[2], results_long[3], results_long[4], results_long[5], results_long[6], results_long[7], results_long[8], results_long[9],
	)
	fmt.Printf(
		"\nShort Ratings Distribution\n[0-10]: %d\n[11-20]: %d\n[21-30]: %d\n[31-40]: %d\n[41-50]: %d\n[51-60]: %d\n[61-70]: %d\n[71-80]: %d\n[81-90]: %d\n[91-100]: %d\n",
		results_short[0], results_short[1], results_short[2], results_short[3], results_short[4], results_short[5], results_short[6], results_short[7], results_short[8], results_short[9],
	)

	fmt.Printf("\nFailed symbols: %v\n", failed)
	return 0
}
