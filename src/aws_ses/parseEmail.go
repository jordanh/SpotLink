package aws_ses

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/jordanh/SpotLink/src/email"
)

func ParseEmail(sesMessage []byte) (*email.EmailData, error) {
	// Unmarshal the JSON into the struct
	var notification Notification
	err := json.Unmarshal(sesMessage, &notification)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}
	// fmt.Println("Decoded SNS message:")
	// pretty.Println(notification)

	// Decode the base64 content
	decodedContent, err := base64.StdEncoding.DecodeString(notification.Content)
	if err != nil {
		fmt.Println("Error decoding base64 content:", err)
		return nil, err
	}

	// Convert the decoded content to a string (or use as is for further processing)
	contentString := string(decodedContent)
	// fmt.Println("Decoded content:", contentString)

	// Parse the email
	emailData, err := email.ParseEmail(contentString)
	if err != nil {
		return nil, err
	}
	return emailData, err
}
