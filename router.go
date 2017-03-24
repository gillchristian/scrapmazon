package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gillchristian/scrapmazon/apm"
	"github.com/julienschmidt/httprouter"
)

// router creates a http.Router, attaches route handlers and returns it.
func router() *httprouter.Router {
	router := httprouter.New()

	router.GET("/movie/amazon/:amazon_id", amazonMovie)

	return router
}

// responseWithJSON sets the response code and writes a json to it.
func responseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

// errorWithJSON sets the response code and writes an error message to it.
func errorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{\"message\": %q}", message)
}

// amazonMovie handles /movie/amazon/:amazon_id route.
func amazonMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Printf("- %v: %v\n", r.Method, r.URL.Path)
	id := ps.ByName("amazon_id")

	movie, err := apm.FetchMovie(id)
	if err != nil {
		// error while fetching Movie
		fmt.Println(err)
		errorWithJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if movie.IsEmpty() {
		// Movie is empty
		msg := "Could not find a movie for amazon_id: " + id
		errorWithJSON(w, msg, http.StatusNotFound)
		return
	}

	data, err := json.Marshal(movie)
	if err != nil {
		// error marshalling Movie data
		fmt.Println(err)
		errorWithJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success!
	responseWithJSON(w, data, http.StatusOK)
}
