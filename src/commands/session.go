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
	Fields      []string
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
		Fields:      fieldList[0:10], // through snr
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
