package main

import (
	"fmt"
	"rate_limiter/limiter"
	"time"
)

func main() {
	capacity := 2
	rateLimiterStrategy := limiter.NewFixedWindowStrategy(int64(capacity), limiter.FIXED_WINDOW_SIZE)

	key1 := limiter.Key("key1")
	key2 := limiter.Key("key2")
	ak1 := rateLimiterStrategy.Allow((key1))
	ak2 := rateLimiterStrategy.Allow((key2))

	fmt.Println(ak1, ak2)

	ak3 := rateLimiterStrategy.Allow((key1))
	ak4 := rateLimiterStrategy.Allow((key2))

	fmt.Println(ak3, ak4)

	time.Sleep(2 * time.Second)

	ak5 := rateLimiterStrategy.Allow((key1))
	ak6 := rateLimiterStrategy.Allow((key2))

	fmt.Println(ak5, ak6)
	// time.Sleep(2 * time.Second)

}
