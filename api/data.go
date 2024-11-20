package api

import "time"

type Bar struct {
	High      float64
	Low       float64
	Open      float64
	Close     float64
	Volume    int64
	Timestamp time.Time
}

type Source string

const (
	Close Source = "close"
	Open  Source = "open"
	High  Source = "high"
	Low   Source = "low"
)
