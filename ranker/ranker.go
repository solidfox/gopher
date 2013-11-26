package ranker

import (
	"fmt"
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
	Size             int64
	Parents          []*spider.Link
	Children         []*spider.Link
}

func NewRanker(option int) *Ranker {
	return &Ranker{option}
}

func (r *Ranker) Search(query *spider.Page) []ResultPage {
	pageIDs, scores := SearchingResult(query)
	fmt.Printf("\npageID: %v\n", pageIDs)
	fmt.Printf("Score: %v\n", scores)
	db := spider.NewRelationalDB("sqlite.db")
	pages := make([]*spider.Page, len(pageIDs))
	//resultPages := make([]*ResultPage, len(pageIDs))
	for i, id := range pageIDs {
		pages[i] = spider.NewPage()
		pages[i].PageID = id
	}
	db.CompleteThePageInfoOf(pages)
	results := make([]ResultPage, len(pages))
	fmt.Println("Loaded page info")
	db.LoadParentsFor(pages)
	db.LoadChildrenFor(pages)

	for i, page := range pages {
		results[i] = ResultPage{
			Title:            page.Title,
			Url:              page.URL,
			Score:            scores[i],
			ModificationDate: page.Modified,
			Size:             page.Size,
			Parents:          page.Parents,
			Children:         page.Children,
		}
	}

	return results
}
