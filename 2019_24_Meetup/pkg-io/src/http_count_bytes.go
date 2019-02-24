package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Counter struct {
	Bytes int
	Lines int
}

func (r *Counter) Write(b []byte) (n int, err error) {
	r.Bytes += len(b)
	r.Lines += bytes.Count(b, []byte{'\n'})
	return len(b), nil
}

var handler = func(w http.ResponseWriter, r *http.Request) {
	cnt := new(Counter)
	io.Copy(io.MultiWriter(w, os.Stdout, cnt), r.Body)
	fmt.Printf("\n===\nBytes: %d\nLines: %d\n", cnt.Bytes, cnt.Lines)
}

// func_main OMIT
func main() {
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handler)))
}
