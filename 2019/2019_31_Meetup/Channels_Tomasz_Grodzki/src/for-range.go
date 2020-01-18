package main

import (
	"fmt"
	"time"
)

func main() {
	// main_begin OMIT
	c := make(chan time.Time, 10)

	// Fill the channel
	for len(c) < cap(c) {
		c <- time.Now()
	}
	close(c) // no more items

	// Read values from the channel
	for v := range c {
		fmt.Println(v.Format(time.RFC3339Nano))
	}
	// main_end OMIT
}
