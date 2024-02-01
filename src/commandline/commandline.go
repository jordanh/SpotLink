package commandline

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/jordanh/SpotLink/src/commands"
	"github.com/tj/go-naturaldate"
)

type CommandLine struct {
	session      *commands.CommandSession
	commandChan  chan commands.Command
	responseChan chan commands.CommandResponse
	done         chan struct{}
}

func NewCommandLine(done chan struct{}) *CommandLine {
	return &CommandLine{
		session:      commands.NewDefaultCommandSession(),
		commandChan:  make(chan commands.Command),
		responseChan: make(chan commands.CommandResponse),
		done:         done,
	}
}

func (cl *CommandLine) handleCommands() {
	for cmd := range cl.commandChan {
		// fmt.Println(pretty.Println(cmd))
		switch cmd.Name {
		case "byCallsign":
			byCallsignCommand, err := commands.NewByCallsignCommand(cmd, cl.session)
			if err != nil {
				cl.responseChan <- commands.CommandResponse{Result: err.Error()}
			} else {
				var sb strings.Builder
				sb.WriteString(
					fmt.Sprintf(
						"Callsign: %s, From Time: %s, To Time: %s",
						byCallsignCommand.Callsign,
						cl.session.FromTime.Format(time.RFC3339),
						cl.session.ToTime.Format(time.RFC3339),
					),
				)
				sb.WriteString("\n\n")
				result, err := commands.ByCallsign(byCallsignCommand)
				if err != nil {
					sb.WriteString(err.Error())
				} else {
					sb.WriteString(result)
				}
				cl.responseChan <- commands.CommandResponse{Result: sb.String()}
			}
		case "help":
			switch len(cmd.ParsedArgs) {
			case 0:
				cl.responseChan <- commands.CommandResponse{Result: `Available commands:
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
					for _, field := range commands.GetValidFields() {
						sb.WriteString(fmt.Sprintln("\t\t - ", field))
					}
					sb.WriteString(`
interval: Natural-language interval before the current time used when searching
          (default: "1 hour ago"). Examples: "6 hours ago" or "1 month ago"
`)
					sb.WriteString(`
   limit: maximum number of query rows to return, integer (default: 10)
`)
					cl.responseChan <- commands.CommandResponse{Result: sb.String()}

				default:
					cl.responseChan <- commands.CommandResponse{Result: `no additional help available`}
				}
			}
		case "set":
			var sb strings.Builder
			for _, parsedArg := range cmd.ParsedArgs {
				switch parsedArg.Type {
				case commands.Word:
					sb.WriteString(fmt.Sprint("warning: skipping unknown word argument: ", parsedArg.Value))
				case commands.KeyValuePair:
					switch parsedArg.Key {
					case "fields":
						sb.WriteString(fmt.Sprintln("fields:", parsedArg.Value))
						err := cl.session.SetFields(strings.Split(parsedArg.Value, ","))
						if err != nil {
							sb.WriteString(fmt.Sprint("error: unable to set fields ", parsedArg.Value))
						}
					case "limit":
						var err error
						cl.session.Limit, err = strconv.Atoi(parsedArg.Value)
						if err != nil {
							sb.WriteString(fmt.Sprint("error: unable to set limit using value ", parsedArg.Value))
						} else {
							sb.WriteString(fmt.Sprintln("limit:", parsedArg.Value))
						}
					case "interval":
						sb.WriteString(fmt.Sprintln("interval:", parsedArg.Value))
						now := time.Now()
						newDate, err := naturaldate.Parse(
							parsedArg.Value,
							time.Now(),
							naturaldate.WithDirection(naturaldate.Past),
						)
						if err == nil {
							cl.session.FromTime, cl.session.ToTime = newDate, now
						} else {
							sb.WriteString(fmt.Sprint("error: unable to parse time interval ", parsedArg.Value))
						}
					case "min_snr":
						var err error
						cl.session.MinSnr, err = strconv.Atoi(parsedArg.Value)
						if err != nil {
							sb.WriteString(fmt.Sprint("error: unable to set min_snr using value ", parsedArg.Value))
						} else {
							sb.WriteString(fmt.Sprintln("min_snr:", parsedArg.Value))
						}
					}
				}
			}
			cl.responseChan <- commands.CommandResponse{Result: sb.String()}
		case "show":
			cl.responseChan <- commands.CommandResponse{Result: cl.session.String()}
		case "q", "quit":
			close(cl.commandChan)
			cl.responseChan <- commands.CommandResponse{Result: "Exiting"}
			cl.done <- struct{}{}
		default:
			cl.responseChan <- commands.CommandResponse{Result: "Error: Unknown command"}
		}
	}
}

func (cl *CommandLine) CommandLineLoop() {
	// Initialize readline
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	// Configure tab completion
	rl.Config.AutoComplete = readline.NewPrefixCompleter(
		readline.PcItem("help"),
		readline.PcItem("byCallsign"),
		readline.PcItem("quit"),
		readline.PcItem("set",
			readline.PcItem("fields"),
			readline.PcItem("interval"),
			readline.PcItem("min_snr"),
			readline.PcItem("limit"),
		),
		readline.PcItem("show"),
	)

	go cl.handleCommands()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		cmd, err := parseCommand(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		cl.commandChan <- cmd

		result := <-cl.responseChan
		fmt.Println(result.Result)

		if cmd.Name == "quit" {
			break
		}
	}
}

func parseCommand(line string) (commands.Command, error) {
	// Regular expression to match key-value pairs and words
	// The key is captured before '=', and the value is captured within quotes
	re := regexp.MustCompile(`(\S+)=("[^"]*"|\S+)|"((?:\\.|[^"\\])*)"|\S+`)
	matches := re.FindAllStringSubmatch(line, -1)

	if len(matches) == 0 {
		return commands.Command{}, fmt.Errorf("invalid command")
	}

	cmd := commands.Command{Name: matches[0][0]}
	for i, match := range matches {
		if len(match[1]) > 0 && len(match[2]) > 0 {
			// This is a key-value pair
			key := match[1]
			value := match[2]

			// Remove quotes if present and unescape any escaped characters
			value = strings.Trim(value, `"`)
			value = strings.Replace(value, `\"`, `"`, -1)
			value = strings.Replace(value, `\\`, `\`, -1)

			cmd.ParsedArgs = append(cmd.ParsedArgs, commands.ParsedArg{Type: commands.KeyValuePair, Key: key, Value: value})
		} else if len(match[0]) > 0 {
			// This is a simple word
			if i != 0 { // Avoid adding the command name as an argument
				value := match[0]
				value = strings.Trim(value, `"`)
				value = strings.Replace(value, `\"`, `"`, -1)
				value = strings.Replace(value, `\\`, `\`, -1)
				cmd.ParsedArgs = append(cmd.ParsedArgs, commands.ParsedArg{Type: commands.Word, Value: value})
			}
		}
	}

	return cmd, nil
}
