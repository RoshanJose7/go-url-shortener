package models

type URL struct {
	Url       string `json:"url"`
	Shortened string `json:"shortened"`
	ExpiresAt string `json:"expires_at"`
}
