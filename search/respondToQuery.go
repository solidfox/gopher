package search

import (
	"gopher/ranker"
	"gopher/spider"
	"io"
	"strings"
)

func respondToQuery(w io.Writer, q string) {
	page := spider.NewPage()
	// TODO Phrases
	page.AddText(q)
	r := ranker.NewRanker()

}
