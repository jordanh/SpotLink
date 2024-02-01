package commands

import (
	"fmt"
	"strings"
	"time"
)

type CommandSession struct {
	FromTime    time.Time
	ToTime      time.Time
	IntervalStr string
	Limit       int
	Fields      []string
	MinSnr      int
}

var fieldList = []string{
	"time", "rx_sign", "rx_loc", "rx_lat", "rx_lon",
	"distance", "azimuth", "frequency", "snr",
	"rx_azimuth", "tx_sign", "tx_loc", "tx_lat", "tx_lon",
	"band", "code", "distance", "drift", "power",
	"id", "version",
}

func NewDefaultCommandSession() *CommandSession {
	fromTime, toTime := timesFromInterval(time.Hour, nil)
	cs := CommandSession{
		FromTime:    fromTime,
		ToTime:      toTime,
		IntervalStr: "1 hour ago",
		Limit:       10,
		Fields:      fieldList[0:10], // through snr
		MinSnr:      0,
	}
	return &cs
}

func timesFromInterval(interval time.Duration, endpoint *time.Time) (time.Time, time.Time) {
	toTime := time.Now()
	if endpoint != nil {
		toTime = *endpoint
	}
	fromTime := toTime.Add(-interval)

	return fromTime, toTime
}

func GetValidFields() []string { return fieldList }

func (cs *CommandSession) SetFields(fields []string) error {
	valid := true
	fieldExists := make(map[string]bool)
	nonExistantFields := make([]string, 0)
	for _, field := range fieldList {
		fieldExists[field] = true
	}
	for _, field := range fields {
		if !fieldExists[field] {
			valid = false
			nonExistantFields = append(nonExistantFields, field)
		}
	}

	if valid {
		cs.Fields = fields
		return nil
	} else {
		return fmt.Errorf("invalid fields %s", strings.Join(nonExistantFields, ", "))
	}

}

func (cs *CommandSession) String() string {
	return fmt.Sprintf(`
Session parameters:
	fields: %s
	interval: %s (from %s to %s)
	limit: %d
	min_snr: %d
`,
		strings.Join(cs.Fields, ", "),
		cs.IntervalStr,
		cs.FromTime.Format(time.RFC3339),
		cs.ToTime.Format(time.RFC3339),
		cs.Limit,
		cs.MinSnr,
	)
}
