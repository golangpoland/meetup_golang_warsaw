package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	// main_begin OMIT
	pr, pw := io.Pipe()
	defer pr.Close()

	// Write stream of data into pipe in a separate routine
	go func() {
		err := json.NewEncoder(pw).Encode(struct {
			Event   string
			Edition int
		}{
			"Golang Meetup Warsaw", 24,
		})

		pw.CloseWithError(err)
	}()

	// Read from pipe
	_, err := http.Post("http://localhost:8080", "application/json", pr)
	// main_end OMIT

	if err != nil {
		log.Fatal(err)
	}
}
