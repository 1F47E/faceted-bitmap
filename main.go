package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	numHotels = 1_000_000
)

// features via bitwise const
const (
	Feature1 uint64 = 1 << iota
	Feature2
	Feature3
	Feature4
)

// same as
// const (
// 	Feature11 uint64 = 1 << 0 // 0x0000000000000001 = 1
// 	Feature22 uint64 = 1 << 1 // 0x0000000000000010 = 2
// 	Feature33 uint64 = 1 << 2 // 0x0000000000000100 = 4
// 	Feature44 uint64 = 1 << 3 // 0x0000000000001000 = 8
// )

func main() {

	// generate batch of hotels with 64 random bits as features
	hotels := make([]uint64, numHotels)
	for i := range hotels {
		hotels[i] = rand.Uint64()
	}

	// mask
	// 0000000000000000000000000000000000000000000000000000000000000101
	mask := Feature1 | Feature3

	fmt.Printf("Mask: %064b\n", mask)

	// parallel search with bitmap mask
	start := time.Now()
	numCPU := runtime.NumCPU()
	chunkSize := numHotels / numCPU
	var wg sync.WaitGroup
	var matchCount int64

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			cnt := int64(0)
			for j := start; j < end; j++ {
				// BITMAP AND
				if hotels[j]&mask == mask {
					cnt++
				}
			}
			atomic.AddInt64(&matchCount, cnt)
		}(i*chunkSize, (i+1)*chunkSize)
	}

	wg.Wait()

	duration := time.Since(start)

	fmt.Printf("total: %d\n", numHotels)
	fmt.Printf("matches: %d out of %d\n", matchCount, numHotels)
	fmt.Printf("took: %v\n", duration)
}
