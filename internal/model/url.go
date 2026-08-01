package model

// ShortenURLRequest is the request payload for shortening a URL.
type ShortenURLRequest struct {
	URL string `json:"url" binding:"required"`
	Exp int    `json:"exp"`
}

// ShortenURLResponse is the response payload for shortening a URL.
type ShortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
