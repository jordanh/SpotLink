package aws_lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kr/pretty"
)

// HandleRequest is your Lambda function handler
func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) (string, error) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("Received SNS message: %s\n", snsRecord.Message)

		// Pretty print the JSON message
		var msg interface{}
		if err := json.Unmarshal([]byte(snsRecord.Message), &msg); err != nil {
			log.Printf("Error unmarshalling SNS message: %v", err)
			return "", err
		}
		pretty.Println(msg)
	}

	return "Successfully processed SNS event", nil
}

func StartLambda() {
	lambda.Start(HandleRequest)
}
