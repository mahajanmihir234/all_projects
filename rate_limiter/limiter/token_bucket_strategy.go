package limiter

import (
	"time"
)

var REFILL_RATE = 0.5

type TokenBucketStrategy struct {
	capacity int
	buckets  map[Key]*Bucket
}

type Bucket struct {
	capacity      int
	tokens        int
	ratePerSecond float64
	lastTime      time.Time
}

func (b *Bucket) refill() {
	now := time.Now()
	elapsed := now.Sub(b.lastTime)
	b.lastTime = now
	tokensToAdd := int(elapsed.Seconds() * b.ratePerSecond)
	if tokensToAdd > 0 {
		b.tokens = min(b.capacity, b.tokens+tokensToAdd)
	}
}

func (b *Bucket) HasToken() bool {
	b.refill()
	if b.tokens <= 0 {
		return false
	}
	b.tokens -= 1
	return true
}

func NewBucket(capacity int, ratePerSecond float64) *Bucket {
	return &Bucket{
		capacity:      capacity,
		tokens:        capacity,
		ratePerSecond: ratePerSecond,
		lastTime:      time.Now(),
	}
}

func (s TokenBucketStrategy) Allow(key Key) bool {
	if _, ok := s.buckets[key]; !ok {
		s.buckets[key] = NewBucket(s.capacity, REFILL_RATE)
	}
	bucket := s.buckets[key]
	return bucket.HasToken()
}

func NewTokenBucketStrategy(capacity int) TokenBucketStrategy {
	return TokenBucketStrategy{
		capacity: capacity,
		buckets:  map[Key]*Bucket{},
	}
}
