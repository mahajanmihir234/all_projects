package shortener

import (
	"errors"
	"time"
)

type URLShortenerService struct {
	domain                string
	urlRepository         URLRepository
	keyGenerationStrategy KeyGenerationStrategy
	maxRetries            int
}

// def shorten(self, long_url: str) -> str:
// # Check if we've already shortened this URL
// existing_key = self.url_repository.find_key_by_long_url(long_url)
// if existing_key is not None:
// 	return self.domain + existing_key

// # Generate a new key, handling potential collisions
// short_key = self._generate_unique_key()

// shortened_url = ShortenedURL.Builder(long_url, short_key).build()
// self.url_repository.save(shortened_url)

// self._notify_observers(EventType.URL_CREATED, shortened_url)

// return self.domain + short_key

func (s URLShortenerService) Shorten(longURL string) (*string, error) {
	shortKey := s.urlRepository.FindShortURLByLongURL(longURL)
	if shortKey != nil {
		key := s.domain + *shortKey
		return &key, nil
	}
	shortKey, err := s.generateUniqueKey()
	if err != nil {
		return nil, err
	}
	shortenedURL := ShortenedURL{
		shortURL:     *shortKey,
		longURL:      longURL,
		creationTime: time.Now(),
	}
	s.urlRepository.Save(shortenedURL)
	key := s.domain + *shortKey
	return &key, nil
}

func (s URLShortenerService) generateUniqueKey() (*string, error) {
	for range s.maxRetries {
		potentialKey := s.keyGenerationStrategy.GenerateKey(s.urlRepository.NextID())
		if !s.urlRepository.ExistsByShortURL(potentialKey) {
			return &potentialKey, nil
		}
	}
	return nil, errors.New("Failed to generate unique key")
}

// for _ in range(self.MAX_RETRIES):
//             # The ID is passed but may be ignored by some strategies (like random)
//             potential_key = self.key_generation_strategy.generate_key(self.url_repository.get_next_id())
//             if not self.url_repository.exists_by_key(potential_key):
//                 return potential_key  # Found a unique key

//         # If we reach here, we failed to generate a unique key after several attempts.
//         raise RuntimeError(f"Failed to generate a unique short key after {self.MAX_RETRIES} attempts.")
