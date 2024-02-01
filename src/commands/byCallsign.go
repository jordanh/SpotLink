package commands

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/jordanh/SpotLink/src/wspr_live"
)

type ByCallsignCommand struct {
	Callsign string
	FromTime time.Time
	ToTime   time.Time
	Limit    int
	Fields   []string
	MinSnr   int
}

func NewByCallsignCommand(command Command, session *CommandSession) (*ByCallsignCommand, error) {
	if len(command.ParsedArgs) != 1 && command.ParsedArgs[0].Type != Word {
		return nil, fmt.Errorf("invalid number of arguments for byCallsign")
	}
	callsign := strings.ToUpper(command.ParsedArgs[0].Value)
	if !ValidCallsign(callsign) {
		return nil, fmt.Errorf("invalid callsign")
	}

	return &ByCallsignCommand{
		Callsign: callsign,
		FromTime: session.FromTime,
		ToTime:   session.ToTime,
		Limit:    session.Limit,
		Fields:   session.Fields,
		MinSnr:   session.MinSnr,
	}, nil
}

func ByCallsign(byCallsignCommand *ByCallsignCommand) (string, error) {
	options := wspr_live.QueryByCallsignOptions{
		FromTime: byCallsignCommand.FromTime,
		ToTime:   byCallsignCommand.ToTime,
		Limit:    byCallsignCommand.Limit,
		MinSnr:   byCallsignCommand.MinSnr,
	}
	queryResponse, err := wspr_live.QueryByCallsign(
		byCallsignCommand.Callsign,
		options,
	)

	if err != nil {
		return "", err
	}

	df := queryResponse.ToDataFrame().Select(byCallsignCommand.Fields)
	if df.Err != nil {
		return "", df.Err
	}

	var buffer bytes.Buffer
	if err := df.WriteCSV(&buffer); err != nil {
		return "", err
	}

	return prettyFmtCSV(buffer.String())
}

func prettyFmtCSV(csvString string) (string, error) {
	reader := csv.NewReader(strings.NewReader(csvString))
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("error reading CSV: %w", err)
	}

	// Find the maximum width of each column for alignment
	maxWidth := make([]int, len(records[0]))
	for _, record := range records {
		for i, field := range record {
			if len(field) > maxWidth[i] {
				maxWidth[i] = len(field)
			}
		}
	}

	var builder strings.Builder

	// Build the formatted CSV string
	for _, record := range records {
		for i, field := range record {
			if i > 0 {
				builder.WriteString(", ")
			}
			fmt.Fprintf(&builder, "%-*s", maxWidth[i], field)
		}
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
