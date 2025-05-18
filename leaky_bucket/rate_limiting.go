package main

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity int
	leakRate time.Duration
	tokens   int
	lastLeak time.Time
	mu       sync.Mutex
}

func NewLeakyBucket(capacity int, leakyRate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakyRate,
		tokens:   capacity,
		lastLeak: time.Now(),
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	fmt.Println("Now: ", now.String())

	elapsedTime := now.Sub(lb.lastLeak)
	fmt.Println("Elapsed time: ", elapsedTime.String())
	fmt.Println("LeakyRate: ", lb.leakRate.String())

	tokensToAdd := int(elapsedTime / lb.leakRate)
	fmt.Println("TokensToAdd: ", tokensToAdd)
	lb.tokens += tokensToAdd

	if lb.tokens > lb.capacity {
		lb.tokens = lb.capacity
	}

	lb.lastLeak = lb.lastLeak.Add(time.Duration(tokensToAdd) * lb.leakRate)

	if lb.tokens > 0 {
		lb.tokens--
		return true
	}

	return false
}

func main() {
	leakyBucket := NewLeakyBucket(5, time.Millisecond*500)

	for range 10 {
		if leakyBucket.Allow() {
			fmt.Println("Request accepted")
		} else {
			fmt.Println("Request denied")
		}
		time.Sleep(time.Second)
	}
}
