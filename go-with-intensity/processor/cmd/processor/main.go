package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	client := NewAuthenticationClient(
		http.DefaultClient, "https", "us-street.api.smartystreets.com",
		"API_KEY", "API_TOKEN")

	pipeline := NewPipeline(os.Stdin, os.Stdout, client, 8)

	if err := pipeline.Process(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
