package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/jordanh/SpotLink/src/wspr_live"
)

type ByCallsignCommand struct {
	Callsign string
	FromTime time.Time
	ToTime   time.Time
}

func NewByCallsignCommand(command Command) (*ByCallsignCommand, error) {
	if len(command.Args) < 1 || len(command.Args) > 3 {
		return nil, fmt.Errorf("Invalid number of arguments for byCallsign")
	}
	callsign := strings.ToUpper(command.Args[0])
	if !ValidCallsign(callsign) {
		return nil, fmt.Errorf("Invalid callsign")
	}

	now := time.Now()

	// Default from time is one hour ago
	fromTime := now.Add(-1 * time.Hour)
	if len(command.Args) >= 2 && ValidISO8601(command.Args[1]) {
		fromTime, _ = time.Parse(time.RFC3339, command.Args[1])
	} else if len(command.Args) >= 2 {
		return nil, fmt.Errorf("Invalid from time format")
	}

	// Default to time is now
	toTime := now
	if len(command.Args) == 3 && ValidISO8601(command.Args[2]) {
		toTime, _ = time.Parse(time.RFC3339, command.Args[2])
	} else if len(command.Args) == 3 {
		return nil, fmt.Errorf("Invalid to time format")
	}

	return &ByCallsignCommand{
		Callsign: callsign,
		FromTime: fromTime,
		ToTime:   toTime,
	}, nil
}

func ByCallsign(byCallsignCommand ByCallsignCommand) (wspr_live.ApiResponse, error) {
	queryResponse, err := wspr_live.QueryByCallsign(
		byCallsignCommand.Callsign,
		byCallsignCommand.FromTime,
		byCallsignCommand.ToTime,
	)
	if err != nil {
		return queryResponse, err
	}
	return queryResponse, nil
}
