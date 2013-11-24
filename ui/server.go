// Daniel Schlaug 2013-11-23

// This program provides a web interface to the Gopher search engine
//
// It serves on http://localhost.local:8081
package main

import (
	"encoding/json"
	"gopher/search"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type message struct {
	Query string
}

func init() {
	log.SetPrefix("Gopher server: ")
}

func main() {
	port := ":8081"
	s := &http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(ServeSearchEngine),
	}
	go func() {
		log.Printf("starting server at localhost%s", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	select {} // Block forever.
}

func ServeSearchEngine(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	log.Println("Serving " + url)

	switch url {
	case "/api":
		respondToApiCall(w, r.Body)
	default:
		http.ServeFile(w, r, "."+url)
	}

}

func respondToApiCall(w io.Writer, r io.ReadCloser) {
	var mess *message

	indata, _ := ioutil.ReadAll(r)

	err := json.Unmarshal(indata, &mess)
	if err != nil {
		log.Println("error:", err)
	}

	search.RespondToQuery(w, mess.Query)
}
