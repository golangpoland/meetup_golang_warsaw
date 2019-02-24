package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.MultiWriter(w, os.Stdout), r.Body)
	}

	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handler)))
}
