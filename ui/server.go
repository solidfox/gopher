// Daniel Schlaug 2013-11-23

// This program provides a web interface to the Gopher search engine
//
// It serves on http://localhost.local:8081
package ui

import (
	"fmt"
	"gopher/search"
	"log"
	"net/http"
	"strconv"
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
	db := moneywire.NewConnection()
	defer db.Close()

	url := r.URL.Path
	log.Println("Serving " + url)

	fmt.Println(request)

	switch url {
	case "/search":
		db.WriteKiosks(request, w)
		query := r.FormValue("query")
		search.respondToQuery(w, query)
	default:
		http.ServeFile(w, r, "client"+url)
	}

}
