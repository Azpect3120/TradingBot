package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type MetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	Interval      string `json:"4. Interval"`
	OutputSize    string `json:"5. Output Size"`
	TimeZone      string `json:"6. Time Zone"`
}

type TimeSeriesData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type Candle struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
	Time   time.Time
}

type Chart struct {
	MetaData MetaData
	Candles  []Candle
}

type apiResponse struct {
	MetaData   MetaData                  `json:"Meta Data"`
	TimeSeries map[string]TimeSeriesData `json:"Time Series (30min)"`
}

// Generate the chart data for a given ticker. This function
// will return a pointer to a Chart struct or an error if
// the API request fails. The caller is responsible for
// checking the error. The data will be cleaned and converted
// from a map to an array for ease of use later.
func GetTickerChart(ticker string) (*Chart, error) {
	var url string = GenerateURL(ticker)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := parseApiResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", response)

	return cleanResponse(response), nil
}

// Parse the data from the Alpha Vantage API into the
// structures defined above. This function does not
// clean any data, just parses it. Any error will be
// returned to the caller.
func parseApiResponse(r io.Reader) (*apiResponse, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	var response apiResponse
	err = json.Unmarshal(body, &response)
	return &response, err
}

// Cleans the response from the API and returns a pointer
// to a better formatted Chart struct. This function will
// clean the data from the map and turn it into an array.
// Furthermore, the TimeSeriesData struct will be converted
// into a Candle struct for easier use later.
func cleanResponse(r *apiResponse) *Chart {
	var chart Chart

	chart.MetaData = r.MetaData // No change here

	// They cannot be stored as time.Time directly.
	// This is because the these need to be used as keys
	// for the map later.
	var timestamps []string
	for timestamp := range r.TimeSeries {
		timestamps = append(timestamps, timestamp)
	}

	// Sort using the time.Time values
	sort.Slice(timestamps, func(i, j int) bool {
		t1, _ := time.Parse("2006-01-02 15:04:05", timestamps[i])
		t2, _ := time.Parse("2006-01-02 15:04:05", timestamps[j])
		return t1.Before(t2)
	})

	// Create blank candles arrays
	chart.Candles = make([]Candle, len(timestamps))

	for _, timestamp := range timestamps {
		data := r.TimeSeries[timestamp]

		openData, _ := strconv.ParseFloat(data.Open, 64)
		highData, _ := strconv.ParseFloat(data.High, 64)
		lowData, _ := strconv.ParseFloat(data.Low, 64)
		closeData, _ := strconv.ParseFloat(data.Close, 64)
		volumeData, _ := strconv.ParseInt(data.Volume, 10, 64)
		timeData, _ := time.Parse("2006-01-02 15:04:05", timestamp)

		var candle Candle = Candle{
			Open:   openData,
			High:   highData,
			Low:    lowData,
			Close:  closeData,
			Volume: volumeData,
			Time:   timeData,
		}

		fmt.Printf("%+v\n", candle)

		chart.Candles = append(chart.Candles, candle)
	}

	return &chart
}
