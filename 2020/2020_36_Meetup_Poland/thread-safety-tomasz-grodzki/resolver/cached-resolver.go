package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var Domains = []string{
	"gmail.com",
	"google.com",
	"youtube.com",
	"facebook.com",
	"yahoo.com",
	"wikipedia.org",
	"amazon.com",
	"live.com",
	"reddit.com",
	"netflix.com",
	"zoom.us",
	"blogspot.com",
	"office.com",
	"microsoft.com",
}

type Resolver struct {
	// Addr specificy DNS server address (IP:port)
	Addr string

	// once will be used for lazy initialization
	once sync.Once

	// resolver is underlying DNS resolver (thread-safe)
	resolver *net.Resolver

	// lock and cache
	lock  sync.RWMutex
	cache map[string][]string

	// cacheHit counts number of lookups resolved via cache
	cacheHits int64
}

// checkCache returns cached result for host (along with true);
// returns (nil, false) if result is not in the cache.
func (r *Resolver) checkCache(host string) ([]string, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	addr, ok := r.cache[host]
	return addr, ok
}

// init initiates the resolver; should be called once.
func (r *Resolver) init() {
	r.cache = make(map[string][]string)

	addr := r.Addr
	if addr == "" {
		addr = "8.8.8.8:53"
	}

	r.resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("udp", addr)
		},
	}
}

// LookupHosts resolves host addresses via DNS.
// It is safe to call LookupHost from multiple goroutines.
func (r *Resolver) LookupHost(host string) ([]string, error) {
	// Initialize resolver on the first call
	r.once.Do(r.init)

	// Check cache first
	if addr, ok := r.checkCache(host); ok {
		// atomically increment cacheHits counter
		atomic.AddInt64(&r.cacheHits, 1)
		return addr, nil
	}

	// Lookup host
	addr, err := r.resolver.LookupHost(context.Background(), host)
	if err != nil {
		return addr, err
	}

	// Save in cache
	r.lock.Lock()
	r.cache[host] = addr
	r.lock.Unlock()

	return addr, nil
}

// CacheHits returns number of calls to LookupHosts which has returned cached values.
// It is safe to call CacheHits from multiple goroutines.
func (r *Resolver) CacheHits() int64 {
	return atomic.LoadInt64(&r.cacheHits)
}

func main() {
	log.SetOutput(os.Stdout)
	start := time.Now()

	// create a resolver using Quad1 server
	resolver := &Resolver{
		Addr: "1.1.1.1:53",
	}

	// wait group will allow us to waint for goroutines to finish
	var wg sync.WaitGroup

	// number of performed lookups (network or cache)
	var lookups int64

	// spawn goroutines performing lookups, but using a single (thread-safe) resolver
	for n := 0; n < 10; n++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for m := 0; m < 100; m++ {
				domain := Domains[rand.Intn(len(Domains))]

				addrs, err := resolver.LookupHost(domain)
				log.Println("DNS:", domain, addrs, err)

				// atomically increment lookup counter
				atomic.AddInt64(&lookups, 1)
			}
		}()
	}

	// spawn goroutine showing live stats
	go func() {
		for {
			// load counters atomically
			log.Printf("Stats: %d lookups, %d cached\n",
				atomic.LoadInt64(&lookups),
				resolver.CacheHits(),
			)
			time.Sleep(time.Millisecond)
		}
	}()

	// wait for all the goroutines to finish
	wg.Wait()

	log.Printf("Done in %s (%d lookups, %d cached)\n",
		time.Since(start),
		lookups,
		resolver.CacheHits(),
	)
}
