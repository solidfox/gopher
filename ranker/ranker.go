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
	pageIDs := SearchingResult(query, r.option)
	db := spider.NewRelationalDB("sqlite.db")
	pages := make([]*spider.Page, len(pageIDs))
	//resultPages := make([]*ResultPage, len(pageIDs))
	for i, id := range pageIDs {
		pages[i] = spider.NewPage()
		pages[i].PageID = id
	}
	db.CompleteThePageInfoOf(pages)
	return []*ResultPage{&ResultPage{Title: query.Words()[0].Word}}
}
