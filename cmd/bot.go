package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/Azpect3120/TradingBot/internal/util"
)

func main() {
	util.GetTickerChart("AAPL")
}

// Deprecated
func get_tickers() []string {
	f, err := os.Open("./data/nasdaq_names.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var tickers []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tickers = append(tickers, strings.TrimSpace(scanner.Text()))
	}

	return tickers
}
