package main

import (
	"fmt"
	"sync"
)

func isPrime(n int) bool {
	for m := 2; m*m <= n; m++ {
		if n%m == 0 {
			return false
		}
	}
	return true
}

func main() {
	check := make(chan int)
	primes := make(chan int)

	// Send numbers to check (jobs)
	go func() {
		for n := 3; n < 100; n += 2 {
			check <- n
		}
		close(check)
	}()

	numOfWorkers := 4

	// Run workers
	// workers_begin OMIT
	var wg sync.WaitGroup

	for n := 0; n < numOfWorkers; n++ {
		wg.Add(1) // HL
		go func() {
			defer wg.Done() // HL
			for n := range check {
				if isPrime(n) {
					primes <- n // HL
				}
			}
		}()
	}

	go func() {
		wg.Wait()     // HL
		close(primes) // HL
	}()
	// workers_end OMIT

	// results_begin OMIT
	for n := range primes {
		fmt.Printf("prime: %d\n", n)
	}
	// results_end OMIT
}
