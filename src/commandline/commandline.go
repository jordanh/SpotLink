package commandline

import (
	"fmt"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/jordanh/SpotLink/src/commands"
	"github.com/kr/pretty"
)

type CommandLine struct {
	commandChan  chan commands.Command
	responseChan chan commands.CommandResponse
	done         chan struct{}
}

func NewCommandLine(done chan struct{}) *CommandLine {
	return &CommandLine{
		commandChan:  make(chan commands.Command),
		responseChan: make(chan commands.CommandResponse),
		done:         done,
	}
}

func (cl *CommandLine) handleCommands() {
	for cmd := range cl.commandChan {
		switch cmd.Name {
		case "byCallsign":
			byCallsignCommand, err := commands.NewByCallsignCommand(cmd)
			if err != nil {
				cl.responseChan <- commands.CommandResponse{Result: err.Error()}
			} else {
				var builder strings.Builder
				builder.WriteString(
					fmt.Sprintf(
						"Callsign: %s, From Time: %s, To Time: %s",
						byCallsignCommand.Callsign,
						byCallsignCommand.FromTime.Format(time.RFC3339),
						byCallsignCommand.ToTime.Format(time.RFC3339),
					),
				)
				builder.WriteString("\n")
				queryResponse, err := commands.ByCallsign(*byCallsignCommand)
				if err != nil {
					builder.WriteString(err.Error())
				} else {
					builder.WriteString(pretty.Sprint(queryResponse))
				}
				cl.responseChan <- commands.CommandResponse{Result: builder.String()}
			}
		case "help":
			cl.responseChan <- commands.CommandResponse{Result: `Available commands:
- byCallsign callsign [from_time] [to_time]: Search by callsign with time range
- help: List all commands
- quit: Quit
`}
		case "quit":
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
	)

	go cl.handleCommands()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		cmd := parseCommand(line)
		cl.commandChan <- cmd

		result := <-cl.responseChan
		fmt.Println(result.Result)

		if cmd.Name == "quit" {
			break
		}
	}
}

func parseCommand(line string) commands.Command {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return commands.Command{}
	}
	return commands.Command{Name: parts[0], Args: parts[1:]}
}
