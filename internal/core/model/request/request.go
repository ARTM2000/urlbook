package request

type request struct {
	TrackId string `json:"track_id"`
}

type SubmitUrl struct {
	request
	Url string `json:"url" validate:""`
}
