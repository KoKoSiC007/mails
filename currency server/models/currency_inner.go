package models

type CurrencyInner struct {
	// Currency name
	Name string `json:"name,omitempty"`
	// Maximum value for passed period
	MaxRate float32 `json:"maxRate,omitempty"`
	// Minimum value for passed period
	MinRate float32 `json:"minRate,omitempty"`
	// Average value for passed period
	AvgRate float32 `json:"avgRate,omitempty"`
}
