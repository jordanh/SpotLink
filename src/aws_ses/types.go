package aws_ses

// Define structs corresponding to the JSON structure
type Notification struct {
	NotificationType string  `json:"notificationType"`
	Mail             Mail    `json:"mail"`
	Receipt          Receipt `json:"receipt"`
	Content          string  `json:"content"`
}

type Mail struct {
	Timestamp        string        `json:"timestamp"`
	Source           string        `json:"source"`
	MessageId        string        `json:"messageId"`
	Destination      []string      `json:"destination"`
	HeadersTruncated bool          `json:"headersTruncated"`
	Headers          []Header      `json:"headers"`
	CommonHeaders    CommonHeaders `json:"commonHeaders"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CommonHeaders struct {
	ReturnPath string   `json:"returnPath"`
	From       []string `json:"from"`
	Date       string   `json:"date"`
	To         []string `json:"to"`
	MessageId  string   `json:"messageId"`
	Subject    string   `json:"subject"`
}

type Receipt struct {
	Timestamp            string   `json:"timestamp"`
	ProcessingTimeMillis int      `json:"processingTimeMillis"`
	Recipients           []string `json:"recipients"`
	SpamVerdict          Verdict  `json:"spamVerdict"`
	VirusVerdict         Verdict  `json:"virusVerdict"`
	SpfVerdict           Verdict  `json:"spfVerdict"`
	DkimVerdict          Verdict  `json:"dkimVerdict"`
	DmarcVerdict         Verdict  `json:"dmarcVerdict"`
	Action               Action   `json:"action"`
}

type Verdict struct {
	Status string `json:"status"`
}

type Action struct {
	Type     string `json:"type"`
	TopicArn string `json:"topicArn"`
	Encoding string `json:"encoding"`
}
