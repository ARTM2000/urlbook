package response

import "github.com/artm2000/urlbook/internal/core/entity/message"

type Response[T any] struct {
	TrackId    string          `json:"track_id"`
	Error      bool            `json:"error"`
	Message    message.Message `json:"message"`
	Data       *T              `json:"data"`
	StatusCode int             `json:"status_code"`
}

type SubmitUrl struct {
	ShortUrl string `json:"short_url"`
}

type GetUrlFromPhrase struct {
	DestinationUrl string
}
