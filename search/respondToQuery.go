package search

import (
	"encoding/json"
	"fmt"
	"gopher/ranker"
	"gopher/spider"
	"io"
	"strings"
)

func RespondToQuery(w io.Writer, q string) {
	query := spider.NewPage()

	for _, word := range strings.Fields(q) {
		query.AddText(word)
	}

	r := ranker.NewRanker(0)
	results := r.Search(query)

	encoder := json.NewEncoder(w)
	fmt.Fprintln(w, "[")

	for _, result := range results {
		encoder.Encode(result)
	}

	fmt.Fprintln(w, "]")
}

func RespondToIndex(w io.Writer, q string) {

	r := ranker.NewRanker(0)
	results := r.GetIndex()

	encoder := json.NewEncoder(w)
	fmt.Fprintln(w, "[")

	for _, result := range results {
		encoder.Encode(result)
	}

	fmt.Fprintln(w, "]")
}
