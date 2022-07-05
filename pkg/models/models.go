package models

type Balance struct {
	UserId  string  `json:"userId"`
	UserId2 string  `json:"userId2,omitempty"`
	Balance float64 `json:"balance"`
	Sum     float64 `json:"sum,omitempty"`
}
