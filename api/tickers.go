package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type data struct {
	// Headers ignored `json:"headers"`
	Rows []Ticker `json:"rows"`
}

type status struct {
	Code             int    `json:"rCode"`
	BadCodeMessage   string `json:"bCodeMessage"`
	DeveloperMessage string `json:"developerMessage"`
}

type response struct {
	Data    data   `json:"data"`
	Message string `json:"message"`
	Status  status `json:"status"`
}

type Ticker struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	LastSale  string `json:"lastsale"`
	NetChange string `json:"netchange"`
	PCTChange string `json:"pctchange"`
	Volume    string `json:"volume"`
	MarketCap string `json:"marketcap"`
	Country   string `json:"country"`
	IPOYear   string `json:"ipoyear"`
	Industry  string `json:"industry"`
	Sector    string `json:"sector"`
	URL       string `json:"url"`
}

// Using the Nasdaq API to get a list of each (7000) ticker
// and their respective information. This function will only
// return an error if there is an error on the API side. Nothing
// directly in this function can fail.
func GetTickers() ([]Ticker, error) {
	var tickers []Ticker
	var url string = "https://api.nasdaq.com/api/screener/stocks?tableonly=true&limit=25&download=true"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Add headers to simulate a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r response
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	tickers = r.Data.Rows

	return tickers, nil
}
