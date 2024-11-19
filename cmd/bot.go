package main

import (
	"bufio"
	"os"
	"strings"
)

func main() {
	tickers := get_tickers()

	for _, ticker := range tickers {
		println(ticker)
	}
}

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
