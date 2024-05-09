package request

type request struct {
	TrackId string `json:"track_id"`
}

type SubmitUrl struct {
	request
	Url string `json:"url" validate:"required,http_url"`
}

type RedirectToDestination struct {
	request
	ShortPhrase string `params:"urlRedirect" validate:"required"`
}
