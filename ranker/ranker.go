package ranker

import (
	"fmt"
	"gopher/spider"
	"strconv"
	"strings"
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
	Keywords         []keyword
	Score            float64
	ModificationDate time.Time
	Size             int64
	Parents          []*spider.Link
	Children         []*spider.Link
}

type keyword struct {
	Word string
	Freq int
}

func NewRanker(option int) *Ranker {
	return &Ranker{option}
}

func (r *Ranker) Search(query *spider.Page) []ResultPage {
	pageIDs, scores, top5words := SearchingResult(query)
	fmt.Print("Top 5 words:")
	for _, str := range top5words {
		fmt.Printf("   %v", str)
	}
	// var mykeyword []keyword
	// for _, top5word := range top5words {
	// 	tuples := strings.Split(top5word, ";")
	// 	for _, tuple := range tuples {
	// 		var tempkeyword keyword
	// 		freqWordWithFreq := strings.Fields(tuple)
	// 		if len(freqWordWithFreq) == 0 {
	// 			tempkeyword.Word = ""
	// 			tempkeyword.Freq = 0
	// 		} else {
	// 			tempkeyword.Word = freqWordWithFreq[0]
	// 			integer, _ := strconv.ParseInt(freqWordWithFreq[1], 10, 64)
	// 			tempkeyword.Freq = int(integer)

	// 		}
	// 		mykeyword = append(mykeyword, tempkeyword)
	// 	}
	// 	for _, tempkeyword := range mykeyword {
	// 		fmt.Printf("keyword:%v  Freq:%v  \n", tempkeyword.Word, tempkeyword.Freq)
	// 	}
	// }

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

		var mykeyword []keyword
		tuples := strings.Split(top5words[i], ";")
		for _, tuple := range tuples {
			var tempkeyword keyword
			freqWordWithFreq := strings.Fields(tuple)
			if len(freqWordWithFreq) == 0 {
				tempkeyword.Word = ""
				tempkeyword.Freq = 0
			} else {
				tempkeyword.Word = freqWordWithFreq[0]
				integer, _ := strconv.ParseInt(freqWordWithFreq[1], 10, 64)
				tempkeyword.Freq = int(integer)

			}
			mykeyword = append(mykeyword, tempkeyword)
		}
		results[i] = ResultPage{
			Title:            page.Title,
			Url:              page.URL,
			Keywords:         mykeyword,
			Score:            scores[i],
			ModificationDate: page.Modified,
			Size:             page.Size,
			Parents:          page.Parents,
			Children:         page.Children,
		}
	}

	return results
}
