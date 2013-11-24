package ranker

import (
	"gopher/spider"
	"time"
)

type Ranker struct {
	option int
}

type Link struct {
	title string
	url   string
}

type ResultPage struct {
	Title            string
	Url              string
	Description      string
	Score            float64
	ModificationDate time.Time
	Size             int
	Parents          []Link
	Children         []Link
}

func NewRanker(option int) *Ranker {
	return &Ranker{option}
}

func (r *Ranker) Search(query *spider.Page) []*ResultPage {
	return []*ResultPage{&ResultPage{Title: "goj"}}
}
