package aws_lambda

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jordanh/SpotLink/src/aws_ses"
	"github.com/kr/pretty"
)

// HandleRequest is your Lambda function handler
func handleRequest(ctx context.Context, snsEvent events.SNSEvent) (string, error) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("Received SNS message: %s\n", snsRecord.Message)
		emailData, err := aws_ses.ParseEmail([]byte(snsRecord.Message))
		if err != nil {
			return "", err
		}
		fmt.Println("Parsed email:")
		pretty.Println(emailData)

		result, err := evalMessage(emailData.Message)
		if err != nil {
			fmt.Println("Error evaluating message:", err)
		}
		err = aws_ses.SendEmail(
			ctx,
			emailData.From,
			emailData.To,
			fmt.Sprint("RE: ", emailData.Subject),
			result,
		)
		if err != nil {
			fmt.Println("Error sending reply:", err)
		}
	}

	return "Successfully processed SNS event", nil
}

func StartLambda() {
	lambda.Start(handleRequest)
}

func TestSnsMessage(testJsonPath string) {
	// Read the test JSON file
	testJsonFile, err := os.Open(testJsonPath)
	if err != nil {
		fmt.Println("Error opening test JSON file:", err)
		return
	}
	defer testJsonFile.Close()

	// Read the file into a byte array
	fileContents, err := os.ReadFile(testJsonPath)
	if err != nil {
		fmt.Printf("Failed to read the JSON file '%s': %s\n", testJsonPath, err)
		return
	}

	emailData, err := aws_ses.ParseEmail(fileContents)
	if err != nil {
		return
	}

	// TODO remove this hard coding
	// testMessage := `	set fields=time,frequency,rx_sign,azimuth,distance,snr interval="1 month ago" limit=20
	// byCallsign k0jrh`
	testMessage := " byCallsign k0jrh\r\n"
	emailData.Message = testMessage

	result, err := evalMessage(emailData.Message)
	if err != nil {
		fmt.Println("Error evaluating message:", err)
	}
	fmt.Println("Result:\n\n", result)
}
