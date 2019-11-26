package main

import "fmt"

func main() {
	// main_begin OMIT
	c := make(chan int, 2)
	c <- 10
	c <- 20
	close(c)

	for n := 0; n < 5; n++ {
		v, ok := <-c
		fmt.Println(v, ok)
	}

	close(c) // panic!
	// main_end OMIT
}
