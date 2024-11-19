package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const API_KEY string = "YZF4227C24SR4F9B"

func main() {
	// tickers := get_tickers()

	// for _, ticker := range tickers {
	// 	println(ticker)
	// }

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=rf&interval=30min&outputsize=compact&datatype=json&extended_hours=false&adjusted=false&apikey=%s", API_KEY)
	println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
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
