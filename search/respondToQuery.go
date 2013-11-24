package search

import (
	"encoding/json"
	"fmt"
	"gopher/ranker"
	"gopher/spider"
	"io"
)

func RespondToQuery(w io.Writer, q string) {
	query := spider.NewPage()
	// TODO Phrases
	query.AddText(q)

	r := ranker.NewRanker(0)
	results := r.Search(query)

	encoder := json.NewEncoder(w)
	fmt.Fprintln(w, "[")

	for _, result := range results {
		encoder.Encode(
			ranker.ResultPage{
				Title: result.Title,
			},
		)
	}

	fmt.Fprintln(w, "]")
}
