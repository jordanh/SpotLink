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

// func validISO8601(dateTime string) bool {
// 	_, err := time.Parse(time.RFC3339, dateTime)
// 	return err == nil
// }
