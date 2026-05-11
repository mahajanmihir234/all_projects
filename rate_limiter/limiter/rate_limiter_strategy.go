package limiter

type Key string

type RateLimiterStrategy interface {
	Allow(key Key) bool
}
