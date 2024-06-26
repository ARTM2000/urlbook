package message

type Message string

const (
	INTERNAL_SYSTEM_ERROR  Message = "internal system error"
	INVALID_URL            Message = "invalid url"
	NOT_FOUND_URL          Message = "not found url"
	DUPLICATE_SHORT_PHRASE Message = "duplicate short url phrase"
)

const (
	SHORT_URL_CREATED Message = "short url created"
)
