package currency

type Currency struct {
	Id       uint   `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Enable   bool   `json:"enable,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}
