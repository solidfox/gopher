package search

import (
	"encoding/json"
	"fmt"
	"gopher/ranker"
	"gopher/spider"
	"io"
	// "strings"
)

func RespondToQuery(w io.Writer, q []string) {
	query := spider.NewPage()

	for _, word := range q {
		query.AddQueryWord(word)
	}

	r := ranker.NewRanker(0)
	results := r.Search(query)

	encoder := json.NewEncoder(w)
	fmt.Fprintln(w, "[")

	for i, result := range results {
		encoder.Encode(result)
		if i < len(results)-1 {
			fmt.Fprint(w, ",")
		}
	}

	fmt.Fprintln(w, "]")
}

func RespondToIndex(w io.Writer) {

	db := spider.NewRelationalDB("sqlite.db")

	results := db.GetIndex()

	encoder := json.NewEncoder(w)
	// fmt.Fprintln(w, "[")

	// for _, result := range results {
	encoder.Encode(results)
	// }

	// fmt.Fprintln(w, "]")
}
