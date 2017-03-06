// The scrapmazon program runs an http server that acts as an API
// for Amazon Prime Movies.
package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", router())
	log.Fatal(err)
}
