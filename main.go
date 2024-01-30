package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jordanh/SpotLink/src/wspr_live"
	"github.com/jordanh/SpotLink/src/wspr_live/callsigns"
	"github.com/kr/pretty"
)

func main() {
	// Define command-line flags
	helpFlag := flag.Bool("h", false, "Print help message")
	interactiveFlag := flag.Bool("i", false, "Attach to stdio and present an interactive command line")
	flag.Parse()

	// Handle help flag
	if *helpFlag {
		fmt.Println("Usage of this program:")
		fmt.Println("-h, --help          Print this help message")
		fmt.Println("-i, --interactive   Attach to stdio and present an interactive command line")
		return
	}

	commandChan := make(chan string)
	responseChan := make(chan string)
	done := make(chan struct{})

	go func() {
		for command := range commandChan {
			handleCommand(command, responseChan, done)
		}
	}()

	// Interactive mode
	if *interactiveFlag {
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Println("Enter commands (type 'help' for list of commands):")
			fmt.Print("> ")
			for scanner.Scan() {
				command := scanner.Text()
				commandChan <- command
				fmt.Println(<-responseChan) // Print responses
				fmt.Print("> ")
			}
			close(commandChan)
		}()
	} else {
		// Default commands for non-interactive mode (can be removed or modified)
		go func() {
			commands := []string{
				"help",
				"byCallsign K0JRH 2024-01-25T15:00:00Z 2024-01-28T15:00:00Z",
				"unknownCommand",
			}
			for _, command := range commands {
				commandChan <- command
				fmt.Println(<-responseChan) // Print responses
			}
			close(commandChan)
		}()
	}

	<-done
}

// 	response, err := wspr_live.QueryByCallsign("K0JRH")
// 	if err != nil {
// 		fmt.Println("error parsing JSON: %w", err)
// 	}

// 	pretty.Println(response)

func handleCommand(command string, responseChan chan<- string, done chan struct{}) {
	args := strings.Fields(command)

	if len(args) == 0 {
		responseChan <- "No command entered"
		return
	}

	switch args[0] {
	case "help":
		responseChan <- `Available commands:
- help: List all commands
- byCallsign callsign [from_time] [to_time]: Search by callsign with time range
- quit: Quit
`

	case "byCallsign":
		if len(args) < 2 || len(args) > 4 {
			responseChan <- "Invalid number of arguments for byCallsign"
			return
		}

		callsign := strings.ToUpper(args[1])
		if !callsigns.ValidCallsign(callsign) {
			responseChan <- "Invalid callsign format"
			return
		}

		// Get current time
		now := time.Now()

		// Default from time is one hour ago
		fromTime := now.Add(-1 * time.Hour)
		if len(args) >= 3 && validISO8601(args[2]) {
			fromTime, _ = time.Parse(time.RFC3339, args[2])
		} else if len(args) >= 3 {
			responseChan <- "Invalid from time format"
			return
		}

		// Default to time is now
		toTime := now
		if len(args) == 4 && validISO8601(args[3]) {
			toTime, _ = time.Parse(time.RFC3339, args[3])
		} else if len(args) == 4 {
			responseChan <- "Invalid to time format"
			return
		}

		response := fmt.Sprintf("Callsign: %s, From Time: %s, To Time: %s", callsign, fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339))
		queryResponse, err := wspr_live.QueryByCallsign("K0JRH", fromTime, toTime)
		if err != nil {
			responseChan <- fmt.Sprintf("error parsing JSON: %s", err)
		}
		responseChan <- fmt.Sprintf("%s\n%s\n", response, pretty.Sprint(queryResponse))
	case "quit":
		close(done)
	default:
		responseChan <- "Unknown command: " + args[0]
	}
}

func validISO8601(dateTime string) bool {
	_, err := time.Parse(time.RFC3339, dateTime)
	return err == nil
}
