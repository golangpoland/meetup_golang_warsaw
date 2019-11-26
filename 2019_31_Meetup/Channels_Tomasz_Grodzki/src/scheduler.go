package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// queue_begin OMIT
type queue []int

// Add number to the queue, keep in low-to-high priority order
func (q *queue) Add(n int) {
	*q = append(*q, n)
	sort.Ints(*q)
}

// Top returns item with the highest priority
func (q queue) Top() int {
	if len(q) > 0 {
		return q[len(q)-1]
	}
	return 0
}

// Pop removes item with the highest priority
func (q *queue) Pop() {
	*q = (*q)[:len(*q)-1]
}

// queue_end OMIT

// scheduler_begin OMIT
func scheduler(in <-chan int, out chan<- int) {
	var q queue

	for len(q) > 0 || in != nil {
		var xout chan<- int // = out if something to send
		if len(q) > 0 {
			xout = out
		}

		select {
		case xout <- q.Top():
			q.Pop()
		case n, ok := <-in:
			if ok {
				q.Add(n)
			} else {
				// channel closed, no more items
				in = nil
			}
		}
	}
}

// scheduler_end OMIT

// main_begin OMIT
func main() {
	in := make(chan int)  // incoming tasks
	out := make(chan int) // outgoing tasks (reordered)

	go func() {
		defer close(in)
		for n := 0; n < 10; n++ {
			in <- rand.Intn(1000)
		}
	}()

	go func() {
		defer close(out)
		scheduler(in, out)
	}()

	for n := range out {
		fmt.Println(n)
		time.Sleep(time.Millisecond) // some work
	}
}

// main_end OMIT
