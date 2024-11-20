package api

import (
	"fmt"
	"time"
)

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

// Return the string representation of the bar. The values are
// displayed as [Timestamp] O<Open> H<High> L<Low> C<Close> Vol<Volume>
func (b *Bar) String() string {
	return fmt.Sprintf("[%s] O:%f H:%f L:%f C:%f Vol:%d", b.Timestamp.Format("2006-01-02 15:04:05"), b.Open, b.High, b.Low, b.Close, b.Volume)
}

// Return the string representation of the bar. The values are
// displayed as [Timestamp] O<Open> H<High> L<Low> C<Close> Vol<Volume>
// with the values rounded to the specified precision.
func (b *Bar) StringFixed(precision int) string {
	var format string = fmt.Sprintf("[%%s] O:%%.%df H:%%.%df L:%%.%df C:%%.%df Vol:%%d", precision, precision, precision, precision)
	return fmt.Sprintf(format, b.Timestamp.Format("2006-01-02 15:04:05"), b.Open, b.High, b.Low, b.Close, b.Volume)
}
