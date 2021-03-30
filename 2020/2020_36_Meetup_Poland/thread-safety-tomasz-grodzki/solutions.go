// #1 OMIT

type Counter map[string]int
res := make(chan Counter)

go func() {
	c := make(Counter)
	c["red"] += 1
	c["blue"] += 1
	res <- c
}()

go func() {
	c := make(Counter)
	c["blue"] += 1
	c["black"] += 1
	res <- c
}()

// merge counters into total
total := make(Counter)
for c := range res {
	// ...
}

// ## OMIT

// #2 OMIT

var (
	s int
	lock sync.Mutex
)

go func() {
	lock.Lock()	
	s += 1
	lock.Unlock()
	// ...
}()

go func() {
	lock.Lock()	
	defer lock.Unlock()
	s += 1
}()

// ## OMIT

// #3 OMIT

var s int64

go func() {
	atomic.AddInt64(&s, 1)
}()

go func() {
	atomic.AddInt64(&s, 1)	
}()

// get total (can be loaded at any time)
total := atomic.LoadInt64(&s)

// ## OMIT

// #4 OMIT

var s int64

go func() {
	v1 := atomic.AddInt64(&s, 1)
}()

go func() {
	v2 := atomic.AddInt64(&s, 1)	
	r2 := atomic.LoadInt64(&s)
}()

// ## OMIT