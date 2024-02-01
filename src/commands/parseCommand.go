package commands

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseCommand(line string) (Command, error) {
	// Regular expression to match key-value pairs and words
	// The key is captured before '=', and the value is captured within quotes
	re := regexp.MustCompile(`(\S+)=("[^"]*"|\S+)|"((?:\\.|[^"\\])*)"|\S+`)
	matches := re.FindAllStringSubmatch(line, -1)

	if len(matches) == 0 {
		fmt.Println("line \"", line, "\"")
		return Command{}, fmt.Errorf("invalid command")
	}

	cmd := Command{Name: matches[0][0]}
	for i, match := range matches {
		if len(match[1]) > 0 && len(match[2]) > 0 {
			// This is a key-value pair
			key := match[1]
			value := match[2]

			// Remove quotes if present and unescape any escaped characters
			value = strings.Trim(value, `"`)
			value = strings.Replace(value, `\"`, `"`, -1)
			value = strings.Replace(value, `\\`, `\`, -1)

			cmd.ParsedArgs = append(cmd.ParsedArgs, ParsedArg{Type: KeyValuePair, Key: key, Value: value})
		} else if len(match[0]) > 0 {
			// This is a simple word
			if i != 0 { // Avoid adding the command name as an argument
				value := match[0]
				value = strings.Trim(value, `"`)
				value = strings.Replace(value, `\"`, `"`, -1)
				value = strings.Replace(value, `\\`, `\`, -1)
				cmd.ParsedArgs = append(cmd.ParsedArgs, ParsedArg{Type: Word, Value: value})
			}
		}
	}

	return cmd, nil
}
