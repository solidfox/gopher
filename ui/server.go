// Daniel Schlaug 2013-11-23

// This program provides a web interface to the Gopher search engine
//
// It serves on http://localhost.local:8081
package main

import (
	"gopher/search"
	"log"
	"net/http"
)

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
	case "/search":
		query := r.FormValue("query")
		search.RespondToQuery(w, query)
	default:
		http.ServeFile(w, r, "."+url)
	}

}
