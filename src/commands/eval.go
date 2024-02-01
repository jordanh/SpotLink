package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tj/go-naturaldate"
)

func WithQuiet(quiet bool) EvalOption {
	return func(ctx *EvalContext) {
		ctx.quiet = quiet
	}
}

func NewEvalContext(options ...EvalOption) *EvalContext {
	ctx := &EvalContext{
		session:      NewDefaultCommandSession(),
		commandChan:  make(chan Command),
		responseChan: make(chan CommandResponse),
	}

	// Apply all options to the context
	for _, option := range options {
		option(ctx)
	}

	return ctx
}

func maybeWriteString(cl *EvalContext, sb *strings.Builder, str string) {
	if !cl.quiet {
		sb.WriteString(str)
	}
}

func (cl *EvalContext) GetEvalChans() (chan<- Command, <-chan CommandResponse) {
	go func() {
		for cmd := range cl.commandChan {
			// fmt.Println(pretty.Println(cmd))
			switch cmd.Name {
			case "byCallsign":
				byCallsignCommand, err := NewByCallsignCommand(cmd, cl.session)
				if err != nil {
					cl.responseChan <- CommandResponse{Result: err.Error()}
				} else {
					var sb strings.Builder
					maybeWriteString(cl, &sb,
						fmt.Sprintf(
							"Callsign: %s, From Time: %s, To Time: %s\n\n",
							byCallsignCommand.Callsign,
							cl.session.FromTime.Format(time.RFC3339),
							cl.session.ToTime.Format(time.RFC3339),
						),
					)
					result, err := ByCallsign(byCallsignCommand)
					if err != nil {
						sb.WriteString(err.Error())
					} else {
						sb.WriteString(result)
					}
					cl.responseChan <- CommandResponse{Result: sb.String()}
				}
			case "help":
				switch len(cmd.ParsedArgs) {
				case 0:
					cl.responseChan <- CommandResponse{Result: `Available commands:
	- byCallsign callsign: Search by callsign with time range
	- help [command]: List all commands or help for a specific command
	- set [field1=value1] [field2=field2] ...: Set parameters used by commands
	- show:  Show current parameters used by commands
	- quit: Quit
	`}
				case 1:
					switch cmd.ParsedArgs[0].Value {
					case "set":
						var sb strings.Builder
						sb.WriteString("set [field1=value1] [field2=value2]...\n")
						sb.WriteString(`
	fields: Comma-separated list of fields to return from query in order
			Example: fields=time,freq,rx_sign,rx_loc,snr
	`)

						sb.WriteString("\n")
						sb.WriteString("\tvalid fields:\n")
						for _, field := range GetValidFields() {
							sb.WriteString(fmt.Sprintln("\t\t - ", field))
						}
						sb.WriteString(`
	interval: Natural-language interval before the current time used when searching
			(default: "1 hour ago"). Examples: "6 hours ago" or "1 month ago"
	`)
						sb.WriteString(`
	limit: maximum number of query rows to return, integer (default: 10)
	`)
						cl.responseChan <- CommandResponse{Result: sb.String()}

					default:
						cl.responseChan <- CommandResponse{Result: `no additional help available`}
					}
				}
			case "set":
				var sb strings.Builder
				for _, parsedArg := range cmd.ParsedArgs {
					switch parsedArg.Type {
					case Word:
						sb.WriteString(fmt.Sprint("warning: skipping unknown word argument: ", parsedArg.Value))
					case KeyValuePair:
						switch parsedArg.Key {
						case "fields":
							err := cl.session.SetFields(strings.Split(parsedArg.Value, ","))
							if err != nil {
								sb.WriteString(fmt.Sprint("error: unable to set fields, ", err))
							} else {
								maybeWriteString(cl, &sb, fmt.Sprintln("fields:", parsedArg.Value))
							}
						case "limit":
							var err error
							cl.session.Limit, err = strconv.Atoi(parsedArg.Value)
							if err != nil {
								sb.WriteString(fmt.Sprint("error: unable to set limit using value ", parsedArg.Value))
							} else {
								maybeWriteString(cl, &sb, fmt.Sprintln("limit:", parsedArg.Value))
							}
						case "interval":
							now := time.Now()
							newDate, err := naturaldate.Parse(
								parsedArg.Value,
								time.Now(),
								naturaldate.WithDirection(naturaldate.Past),
							)
							if err != nil {
								sb.WriteString(fmt.Sprint("error: unable to parse time interval ", parsedArg.Value))
							} else {
								cl.session.FromTime, cl.session.ToTime = newDate, now
								maybeWriteString(cl, &sb, fmt.Sprintln("interval:", parsedArg.Value))
							}
						case "min_snr":
							var err error
							cl.session.MinSnr, err = strconv.Atoi(parsedArg.Value)
							if err != nil {
								sb.WriteString(fmt.Sprint("error: unable to set min_snr using value ", parsedArg.Value))
							} else {
								maybeWriteString(cl, &sb, fmt.Sprintln("min_snr:", parsedArg.Value))
							}
						}
					}
				}
				cl.responseChan <- CommandResponse{Result: sb.String()}
			case "show":
				cl.responseChan <- CommandResponse{Result: cl.session.String()}
			case "q", "quit":
				cl.responseChan <- CommandResponse{Result: "Exiting", Final: true}
			default:
				cl.responseChan <- CommandResponse{Result: "Error: Unknown command"}
			}
		}
	}()

	return cl.commandChan, cl.responseChan
}

func (cl *EvalContext) Close() {
	close(cl.commandChan)
	close(cl.responseChan)
}
