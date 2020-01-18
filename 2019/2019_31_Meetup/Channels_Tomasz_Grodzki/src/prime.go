// package_begin OMIT
package main

import (
	"fmt"
	"math"
	"time"
)

func isPrime(n int) bool {
	for m := 2; m <= int(math.Sqrt(float64(n))); m++ {
		if n%m == 0 {
			return false
		}
	}
	return true
}

// ...
// package_end OMIT

// jobs_begin OMIT
func main() {
	start := time.Now()

	check := make(chan int, 100)
	results := make(chan int)

	// Send candidates (all odd numbers in range)
	go func() {
		for n := 3; n < 2e6; n += 2 {
			check <- n // HL
		}
		close(check) // HL
	}()

	// ...
	// jobs_end OMIT

	// workers_begin OMIT
	// Run workers
	numOfWorkers := 4

	for n := 0; n < numOfWorkers; n++ {
		go func() {
			found := 0
			for n := range check { // HL
				if isPrime(n) {
					found++
				}
			}
			results <- found // HL
		}()
	}

	// ...
	// workers_end OMIT

	// results_begin OMIT
	// Collect the results
	found := 0
	for n := 0; n < numOfWorkers; n++ {
		found += <-results // HL
	}

	fmt.Printf("found %d prime numbers in %s\n", found, time.Since(start))
}

// results_end OMIT
