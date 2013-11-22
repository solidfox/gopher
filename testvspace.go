package main

import (
	"fmt"
	"gopher/ranker"
	"gopher/spider"
	"time"
)

func main() {
	start := time.Now()

	db := spider.NewDBM("DBM.db")
	wordN := db.GetWordNumber()
	docN := db.GetDocumentNumber()
	fmt.Printf("Documents stored: %v\n", docN)
	fmt.Printf("Words stored: %v\n", wordN)
	fmt.Printf("Df of wordid=10: %v\n", db.Getdf(10))
	fmt.Printf("Inside stuff 10: ")
	pageIds := db.GetDocIdByWordID(10)
	var words []string
	words = append(words, "societi")
	for _, pageId := range pageIds {
		fmt.Printf("pageId: %v	TF: %v	TFIDF: %v", pageId, int(db.GetTf(10, pageId)), db.GetTfidf(10, pageId))
		fmt.Printf("	CosSim: %v\n", db.CosSim(pageId, words))
	}

	//ranker.PrintHiHi()

	pages2 := db.GetPages2()
	var testingPage spider.Page
	for _, page := range pages2 {
		testingPage = *page
		break
	}
	// for _, word := range testingPage.Words() {
	// 	fmt.Printf("%v", word)
	// }
	db.Close()
	ranker.SearchingResult(testingPage, ranker.TFIDF)
	elapsed := time.Since(start)
	fmt.Printf("Time spent: %v\n", elapsed)

}
