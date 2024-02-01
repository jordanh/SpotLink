package main

import (
	"flag"
	"fmt"

	"github.com/jordanh/SpotLink/src/aws_lambda"
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
		fmt.Println("-h, Print this help message")
		fmt.Println("-i, Attach to stdio and present an interactive command line")
		return
	}

	// Interactive mode
	if *interactiveFlag {
		doInteraciveMode()
	} else {
		aws_lambda.StartLambda()
	}
}
