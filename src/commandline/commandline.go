package commandline

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/jordanh/SpotLink/src/commands"
)

func CommandLineLoop(ctx *commands.EvalContext) error {
	if ctx == nil {
		return fmt.Errorf("ctx is nil")
	}

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

	commandChan, responseChan := ctx.GetEvalChans()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		cmd, err := commands.ParseCommand(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		commandChan <- cmd
		result := <-responseChan
		fmt.Println(result.Result)

		if result.Final {
			break
		}
	}
	return nil
}
