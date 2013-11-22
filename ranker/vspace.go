package ranker

import (
	//"fmt"
	"fmt"
	"gopher/spider"
	//"math"
)

const (
	TFIDF = iota //=0
	Cos          //1
	//add here
)

func Tfidf(documentID int64, wordID int) (score float64) {
	db := spider.NewDBM(spider.DBMname)
	score = db.GetTfidf(wordID, documentID)
	db.Close()
	return
}

func CosSim(docId int64, query []string) (score float64) {
	db := spider.NewDBM(spider.DBMname)
	score = db.CosSim(docId, query)
	db.Close()
	return
}

func QueryProcess(query []string) (docment spider.Page) {
	return
}

func SearchingResult(query spider.Page, option int) (docIDs []int64) {

	switch option {
	default:
		fmt.Errorf("Search with invalid option")
	case TFIDF:
		words := query.Words()
		docIDs := getDocIDsbyWords(words)
		// for _, ids := range words {
		// 	fmt.Printf("\n%v", ids.WordID)
		// }
		docResult := make([]float64, len(docIDs))
		//var docResult [len(docIDs)]float64
		db := spider.NewDBM(spider.DBMname)
		for index, documentID := range docIDs {
			docResult[index] = 0
			for _, word := range words {
				wordID := word.WordID

				docResult[index] += db.GetTfidf(wordID, documentID) //Tfidf(documentID, wordID)

			}
			//fmt.Printf("%v", docResult)
		}
		db.Close()

	case Cos:
		var docResult []float64
		words := query.Words()
		var strArr []string
		for _, word := range words {
			wordContent := word.Word
			strArr = append(strArr, wordContent)
		}
		docIDs := getDocIDsbyWords(words)
		db := spider.NewDBM(spider.DBMname)
		for index, searchDocID := range docIDs {

			docResult[index] = db.CosSim(searchDocID, strArr) //CosSim(searchDocID, strArr)

		}
		db.Close()
		//fmt.Printf("%v", docResult)
	}
	return
}

func removeArrayItemWithIndex(arr []float64, index int) []float64 {
	arr = append(arr[:index], arr[index+1:]...)
	return arr
}

func PrintHiHi() {
	fmt.Printf("HiHi")
}

func getDocIDsbyWords(words []*spider.Word) (DocIDs []int64) {
	db := spider.NewDBM("DBM.db")
	for _, word := range words {
		docIDs := db.GetDocIdByWordID(word.WordID)
		for _, docID := range docIDs {
			if len(DocIDs) == 0 {
				DocIDs = append(DocIDs, docID)
			} else {
				if isExist(DocIDs, docID) == false {
					DocIDs = append(DocIDs, docID)
				}
			}
		}
	}
	db.Close()
	return
}

func isExist(input []int64, newItem int64) (result bool) {
	for _, item := range input {
		if item == newItem {
			return true
		}
	}
	return false
}
