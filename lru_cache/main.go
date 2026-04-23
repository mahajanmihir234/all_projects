package main

import (
	"fmt"
	"lru_cache/cache"
	"sync"
)

func main() {
	c := cache.NewLRUCache(3)
	c.Put(1, 1)
	c.Print()
	c.Put(2, 2)
	c.Print()
	c.Put(3, 3)
	c.Print()

	c.Put(4, 4)
	c.Print()
	c.Put(3, 5)
	c.Print()

	fmt.Println()

	wg := sync.WaitGroup{}

	wg.Add(2)
	var f = func(key interface{}, val interface{}) {
		defer wg.Done()
		c.Put(key, val)
	}
	go f(1, 4)
	go f(1, 3)

	wg.Wait()

	c.Print()
}
