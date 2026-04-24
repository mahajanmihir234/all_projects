package shortener

import "sync"

type URLRepository interface {
	Save(url ShortenedURL) error
	FindByShortURL(shortURL string) *ShortenedURL
	FindShortURLByLongURL(longURL string) *string
	NextID() int
	ExistsByShortURL(shortURL string) bool
}

type InMemoryURLRepository struct {
	shortURLToURLMap     map[string]ShortenedURL
	longURLToShortURLMap map[string]string
	idCounter            int
	mutex                sync.Mutex
}

func (r *InMemoryURLRepository) Save(url ShortenedURL) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.shortURLToURLMap[url.ShortURL()] = url
	r.longURLToShortURLMap[url.LongURL()] = url.ShortURL()
	return nil
}

func (r *InMemoryURLRepository) FindByShortURL(shortURL string) *ShortenedURL {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if url, ok := r.shortURLToURLMap[shortURL]; ok {
		return &url
	}
	return nil
}

func (r *InMemoryURLRepository) FindShortURLByLongURL(longURL string) *string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if shortURL, ok := r.longURLToShortURLMap[longURL]; ok {
		return &shortURL
	}
	return nil
}

func (r *InMemoryURLRepository) NextID() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	value := r.idCounter
	r.idCounter++
	return value
}

func (r *InMemoryURLRepository) ExistsByShortURL(shortURL string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, ok := r.shortURLToURLMap[shortURL]
	return ok
}
