package aws_lambda

import (
	"strings"

	"github.com/jordanh/SpotLink/src/commands"
)

func evalMessage(message string) (string, error) {
	var builder strings.Builder
	ctx := commands.NewEvalContext(commands.WithQuiet(true))
	commandChan, responseChan := ctx.GetEvalChans()
	defer ctx.Close()

	// Split message into trimmed lines, rejecting empty lines
	var lines []string
	for _, line := range strings.Split(message, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			lines = append(lines, line)
		}
	}

	// for each line, parse into Commands:
	for _, line := range lines {
		cmd, err := commands.ParseCommand(line)
		if err != nil {
			// terminate on first error
			return builder.String(), err
		}
		commandChan <- cmd
		result := <-responseChan
		builder.WriteString(result.Result)

		if result.Final {
			break
		}
	}
	return builder.String(), nil
}
