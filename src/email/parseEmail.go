package email

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

func ParseEmail(rawEmail string) (*EmailData, error) {
	msg, err := mail.ReadMessage(strings.NewReader(rawEmail))
	if err != nil {
		return nil, err
	}

	header := msg.Header
	from := header.Get("From")
	to := header.Get("To")
	date := header.Get("Date")
	subject := header.Get("Subject")

	contentType := header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse media type: %w", err)
	}

	var message string
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break // End of the multipart entities
			} else if err != nil {
				return nil, fmt.Errorf("failed to read next part: %w", err)
			}
			partMediaType, _, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
			if err != nil {
				continue // Unable to parse the part's media type, skip it
			}
			if partMediaType == "text/plain" {
				// Prefer text/plain parts
				bodyBytes, err := io.ReadAll(p)
				if err != nil {
					return nil, fmt.Errorf("failed to read text/plain part: %w", err)
				}
				message = string(bodyBytes)
				break // Found the plain text part, no need to look further
			}
		}
	} else {
		// Non-multipart, read the whole body
		bodyBytes, err := io.ReadAll(msg.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
		message = string(bodyBytes)
		// Set the mediaType to text/plain if it's not already set
		if mediaType == "" {
			mediaType = "text/plain"
		}
	}

	emailData := &EmailData{
		From:     from,
		To:       to,
		Date:     date,
		Subject:  subject,
		Message:  message,
		MimeType: mediaType,
	}

	return emailData, nil
}
