package main

import "fmt"

func main() {
	c := make(chan string, 10)
	c <- "konichiwa"

	fmt.Println("cap:", cap(c)) // 10
	fmt.Println("len:", len(c)) // 1
}
