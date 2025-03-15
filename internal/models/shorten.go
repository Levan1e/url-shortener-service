package models

type GetShortenByOriginalRequest struct {
	Url string `json:"url"`
}

type GetShortenByOriginalResponse struct {
	ShortenUrl string `json:"shorten_url"`
}

type GetOriginalByShortenResponse struct {
	Url string `json:"url"`
}
