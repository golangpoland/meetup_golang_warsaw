package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, r.Body)
	}

	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handler)))
}
