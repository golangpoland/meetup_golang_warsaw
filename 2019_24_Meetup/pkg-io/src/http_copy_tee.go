package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, io.TeeReader(r.Body, os.Stdout))
	}

	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handler)))
}
