package main

import (
	"log"
	"time"
)

func main() {
	// semaphore_begin OMIT
	semaphore := make(chan struct{}, 2)

	task := func() bool {
		// non-blocking acquire // HL
		select { // HL
		case semaphore <- struct{}{}: // HL
			defer func() { <-semaphore }() // HL
		default: // HL
			// channel full // HL
			return false // HL
		} // HL

		// some work here
		time.Sleep(time.Second)
		return true
	}
	// semaphore_end OMIT

	// workers_begin OMIT
	for n := 0; n < 10; n++ {
		go func(id int) {
			for {
				if task() {
					log.Printf("[%d]: task done!\n", id)
				} else {
					log.Printf("[%d]: would block, doing something else...\n", id)
					time.Sleep(time.Second / 2)
				}
			}
		}(n)
	}
	// workers_end OMIT

	time.Sleep(time.Minute)
}
