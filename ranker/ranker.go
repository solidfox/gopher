package ranker

import (
	"gopher/spider"
	"time"
)

type Ranker struct {
}

type Link struct {
	title string
	url   string
}

type ResultPage struct {
	Title            string
	Url              string
	Description      string
	Score            float32
	ModificationDate time.Time
	Size             int
	Parents          []Link
	Children         []Link
}

func NewRanker() *Ranker {
	return &Ranker{}
}

func (r *Ranker) Search(query *spider.Page) []*ResultPage {
	dbm := spider.NewDBM(spider.DBMname)
	defer dbm.Close()
	pages, scores := dbm.GetPagesForQuery(query)
	results := make([]*ResultPage, len(pages))
	var rdb *spider.RelationalDB
	rdb = dbm.RDB()
	rdb.LoadParentsFor(pages)
	for i, page := range pages {
		results[i] = &ResultPage{
			Title: page.Title,
			Score: scores[i],
		}
	}
	return results
}
