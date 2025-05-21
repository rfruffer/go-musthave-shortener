package models

type (
	ShortenRequest struct {
		URL string `json:"url"`
	}
	ShortenResponse struct {
		Result string `json:"result"`
	}
)
