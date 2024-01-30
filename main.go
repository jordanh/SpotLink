package main

import (
	"flag"
	"fmt"

	"github.com/jordanh/SpotLink/src/commandline"
)

func doInteraciveMode() {
	done := make(chan struct{})

	go func() {
		cl := commandline.NewCommandLine(done)
		cl.CommandLineLoop()
	}()

	<-done
}

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

	// Interactive mode
	if *interactiveFlag {
		doInteraciveMode()
	} else {
		// TK TODO
	}
}

// func handleCommand(command string, responseChan chan<- string, done chan struct{}) {
// 	args := strings.Fields(command)

// 	if len(args) == 0 {
// 		responseChan <- "No command entered"
// 		return
// 	}

// 	switch args[0] {
// 	case "help":
// 		responseChan <- `Available commands:
// - help: List all commands
// - byCallsign callsign [from_time] [to_time]: Search by callsign with time range
// - quit: Quit
// `

// 	case "byCallsign":
// 		if len(args) < 2 || len(args) > 4 {
// 			responseChan <- "Invalid number of arguments for byCallsign"
// 			return
// 		}

// 		callsign := strings.ToUpper(args[1])
// 		if !callsigns.ValidCallsign(callsign) {
// 			responseChan <- "Invalid callsign format"
// 			return
// 		}

// 		// Get current time
// 		now := time.Now()

// 		// Default from time is one hour ago
// 		fromTime := now.Add(-1 * time.Hour)
// 		if len(args) >= 3 && validISO8601(args[2]) {
// 			fromTime, _ = time.Parse(time.RFC3339, args[2])
// 		} else if len(args) >= 3 {
// 			responseChan <- "Invalid from time format"
// 			return
// 		}

// 		// Default to time is now
// 		toTime := now
// 		if len(args) == 4 && validISO8601(args[3]) {
// 			toTime, _ = time.Parse(time.RFC3339, args[3])
// 		} else if len(args) == 4 {
// 			responseChan <- "Invalid to time format"
// 			return
// 		}

// 		response := fmt.Sprintf("Callsign: %s, From Time: %s, To Time: %s", callsign, fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339))
// 		queryResponse, err := wspr_live.QueryByCallsign("K0JRH", fromTime, toTime)
// 		if err != nil {
// 			responseChan <- fmt.Sprintf("error parsing JSON: %s", err)
// 		}
// 		responseChan <- fmt.Sprintf("%s\n%s\n", response, pretty.Sprint(queryResponse))
// 	case "quit":
// 		close(done)
// 	default:
// 		responseChan <- "Unknown command: " + args[0]
// 	}
// }

// func validISO8601(dateTime string) bool {
// 	_, err := time.Parse(time.RFC3339, dateTime)
// 	return err == nil
// }
