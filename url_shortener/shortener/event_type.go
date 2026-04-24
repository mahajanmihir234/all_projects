package shortener

type EventType string

const (
	URL_CREATED  EventType = "URL_CREATED"
	URL_ACCESSED EventType = "URL_ACCESSED"
)
