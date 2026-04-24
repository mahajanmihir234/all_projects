package shortener

import "time"

type ShortenedURL struct {
	shortURL     string
	longURL      string
	creationTime time.Time
}

func (s ShortenedURL) ShortURL() string {
	return s.shortURL
}

func (s ShortenedURL) LongURL() string {
	return s.longURL
}

func (s ShortenedURL) CreationTime() time.Time {
	return s.creationTime
}
