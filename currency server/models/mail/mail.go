package mail

type Mail struct {
	// ID of email
	Id uint `json:"id,omitempty"`
	// From mail
	From string `json:"from,omitempty"`
	// To mail
	To string `json:"to,omitempty"`
	// Body of mail
	Body string `json:"body,omitempty"`
}
