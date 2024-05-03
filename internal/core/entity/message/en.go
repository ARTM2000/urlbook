package message

type Message string

const (
	INTERNAL_SYSTEM_ERROR Message = "internal system error"
	INVALID_URL Message = "invalid url"
)

const (
	SHORT_URL_CREATED Message = "short url created"
)
