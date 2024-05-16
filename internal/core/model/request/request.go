package request

type request struct {
	TrackId string `json:"track_id"`
}

type SubmitUrl struct {
	request
	Url string `json:"url" validate:"required,http_url"`
}

type SubmitUrlByCustomPhrase struct {
	SubmitUrl
	Phrase string `json:"phrase" validate:"required,alphanum,max=16"`
}

type RedirectToDestination struct {
	request
	ShortPhrase string `params:"urlRedirect" validate:"required"`
}
