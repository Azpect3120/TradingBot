package main

import (
	"fmt"

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

		rating := api.CalculateRating(row, bars)
		if rating.LongScore >= 50 || rating.ShortScore >= 50 {
			fmt.Printf("Rating %s Long: %0.2f Short: %0.2f\n", row, rating.LongScore, rating.ShortScore)
		}
	}
}
