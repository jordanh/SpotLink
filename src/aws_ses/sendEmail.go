package aws_ses

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/kr/pretty"
)

func SendEmail(ctx context.Context, to, from, subject, body string) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return fmt.Errorf("error loading AWS configuration: %v", err)
	}

	svc := ses.NewFromConfig(cfg)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(body),
				},
			},
			Subject: &types.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(from),
	}
	fmt.Println("Will send:")
	pretty.Println(input)

	_, err = svc.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}
