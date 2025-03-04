package mail

type MailBody struct {
	// To mail
	To string `json:"to,omitempty"`
	// Body of mail
	Data string `json:"data,omitempty"`
}
