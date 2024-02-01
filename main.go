package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jordanh/SpotLink/src/aws_lambda"
	"github.com/jordanh/SpotLink/src/commandline"
	"github.com/jordanh/SpotLink/src/commands"
)

func doInteraciveMode() {
	ctx := commands.NewEvalContext()
	defer ctx.Close()
	commandline.CommandLineLoop(ctx)
}

func main() {
	// Define command-line flags
	helpFlag := flag.Bool("h", false, "Print help message")
	interactiveFlag := flag.Bool("i", false, "Attach to stdio and present an interactive command line")
	testJsonPath := flag.String("t", "", "Enable test mode, path to a test SNS Message JSON file")
	flag.Parse()

	// Handle help flag
	if *helpFlag {
		fmt.Println("Usage of this program:")
		fmt.Println("-h, Print this help message")
		fmt.Println("-i, Attach to stdio and present an interactive command line")
		fmt.Println("-t <path>, Path to a JSON file for testing")
		return
	}

	if *testJsonPath != "" {
		if _, err := os.Stat(*testJsonPath); os.IsNotExist(err) {
			fmt.Printf("The JSON file '%s' does not exist.\n", *testJsonPath)
			return
		}
		fmt.Printf("Using JSON file for testing: %s\n", *testJsonPath)
	}

	// Interactive mode
	if *interactiveFlag {
		doInteraciveMode()
	} else if *testJsonPath != "" {
		aws_lambda.TestSnsMessage(*testJsonPath)
	} else {
		aws_lambda.StartLambda()
	}
}
