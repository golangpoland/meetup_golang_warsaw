package main

import (
	"log"
	"time"
)

func main() {
	// main_begin OMIT
	semaphore := make(chan struct{}, 2) // HL

	task := func() {
		semaphore <- struct{}{}        // acquire // HL
		defer func() { <-semaphore }() // release // HL

		// some work here
		time.Sleep(time.Second)
	}

	for n := 0; n < 10; n++ {
		go func(id int) {
			log.Printf("[%d]: starting...\n", id)
			for {
				task()
				log.Printf("[%d]: task done!\n", id)
			}
		}(n)
	}
	// main_end OMIT

	time.Sleep(time.Minute)
}
