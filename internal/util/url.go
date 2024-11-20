package util

import "fmt"

const api_key string = "YZF4227C24SR4F9B"

func GenerateURL(symbol string) string {
	var interval string = "30min"     // 1min, 5min, 15min, 30min, 60min
	var outputsize string = "compact" // compact, full

	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=%s&interval=%s&outputsize=%s&extended_hours=false&adjusted=false&apikey=%s",
		symbol,
		interval,
		outputsize,
		api_key,
	)

	return url
}
