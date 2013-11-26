package ranker

import (
	"fmt"
	"gopher/spider"
	"math"
	"strconv"
	"strings"
	"time"
)

type wordWithTFIDF struct {
	Word  *spider.Word
	TFIDF float64
}

type SPage struct {
	Page   *spider.Page
	myWord []wordWithTFIDF
}

func contain(list []int64, item int64) bool {
	for _, temp := range list {
		if temp == item {
			return true
		}
	}
	return false
}

func SearchingResult(query *spider.Page) (resultPageIDs []int64, resultScores []float64) {
	//query.wo

	//storing page with TFIDF
	resultScores = nil
	resultPageIDs = nil
	db := spider.NewDBM(spider.DBMname)
	invertedindex := db.GetInvertedIndex()
	allStoredPages := db.GetPages2()
	//computeAveDocLen(allStoredPages)
	invertedTable := make(map[int][]int64)
	for index, temp := range invertedindex {
		terms := strings.Split(temp, ";")
		for _, term := range terms {
			docID, _ := strconv.ParseInt(term, 10, 64)
			if contain(invertedTable[index], docID) == false {
				invertedTable[index] = append(invertedTable[index], docID)
			}
		}
	}

	var allPagesWithTFIDF []SPage
	//compute allPagesWithTFIDF
	fmt.Printf("Document number in allstoredpage:%v", len(allStoredPages))
	for _, page := range allStoredPages {
		var tempPage SPage
		tempPage.Page = page
		for _, word := range page.Words() {
			var tempWord wordWithTFIDF
			tempWord.Word = word
			TF := float64(word.TF())
			DF := float64(len(invertedTable[word.WordID]))
			N := float64(len(allStoredPages))
			MaxTF := GetMaxTF(page)
			//fmt.Printf("TF:%v    DF:%v\n", TF, DF)
			if MaxTF <= 0 {
				//fmt.Printf("'%v' is word in the db with 0 MAXTF", word.Word)
				tempWord.TFIDF = 0
			} else if DF <= 0 {
				//fmt.Printf("'%v' is word in the db with 0 DF", word.Word)
				tempWord.TFIDF = 0
			} else {
				tempWord.TFIDF = (TF / MaxTF) * math.Log2(N/DF)
			}
			tempPage.myWord = append(tempPage.myWord, tempWord)
		}

		allPagesWithTFIDF = append(allPagesWithTFIDF, tempPage)
	}
	//fmt.Printf("len of allPagesWithTFIDF:%v	\n", len(allPagesWithTFIDF))
	// for _, pageWithTFIDF := range allPagesWithTFIDF {
	// 	//allPagesWithTFIDF
	// 	words := pageWithTFIDF.myWord
	// 	// fmt.Printf("format")
	// 	for _, word := range words {
	// 		str := word.Word.Word
	// 		TFIDF := word.TFIDF
	// 		fmt.Printf("Str:%v   TFIDF:%v\n", str, TFIDF)
	// 	}
	// }

	//end of storing TFIDF

	//handling querry
	Start := time.Now()
	var scores []float64
	var pageIDs []int64
	QWords := query.Words()
	// QWords[0].Word = "shuten" + " " + "doji"
	// QWords = QWords[:2]
	for _, word := range QWords {
		fmt.Printf("Qword:%v\n", word.Word)
	}
	// for _, page := range allStoredPages {
	// 	str := QWords[0].Word
	// 	if strings.Contains(str, " ") {
	// 		phaseterms := strings.Fields(str)
	// 		temp := GetTFIDFPhased(page, allStoredPages, invertedTable, phaseterms)
	// 		if temp > 0 {
	// 			fmt.Printf("GOOD")
	// 		}
	// 	}

	// }
	for _, page := range allPagesWithTFIDF {
		var dq float64 = 0.0
		var dlen float64 = GetDLen(page)
		var qlen float64 = float64(len(query.Words()))

		for _, querryWord := range QWords {
			//str := querryWord.Word
			words := page.myWord

			str := querryWord.Word
			if strings.Contains(str, " ") {
				phaseterms := strings.Fields(str)
				// Start := time.Now()
				temp := GetTFIDFPhased(page.Page, allStoredPages, invertedTable, phaseterms)

				dq += temp
				if temp != 0 {
					fmt.Printf("temp: %v\n", temp)
					fmt.Printf("Page with phased:%v\n", page.Page.URL)
				}

			}

			//fmt.Printf("phase: %v\n", QWords[0].Word)
			for _, word := range words {
				if querryWord.Word == word.Word.Word {
					dq += word.TFIDF * 1
					//fmt.Print("notPhase")
				}
			}
		}
		scores = append(scores, dq/dlen/qlen)
		pageIDs = append(pageIDs, page.Page.PageID)
	}

	for _, word := range allPagesWithTFIDF[53-1].myWord {
		if word.Word.Word == "demon" {
			fmt.Printf("str: %v", word.Word.Word)
		}

	}
	fmt.Printf("\n%v\n", allPagesWithTFIDF[53-1].Page.URL)
	fmt.Printf("pageIDs:%v    ", pageIDs)
	fmt.Printf("scores:%v\n", scores)

	// for _, word := range query.Words() {
	// 	fmt.Printf("querry words:%v\n", word.Word)
	// }

	MaxNumPageReturn := 50
	for i := 0; i < MaxNumPageReturn; i++ {
		maxIndex := 0
		maxValue := 0.0
		for index, score := range scores {
			if score > maxValue {
				maxIndex = index
				maxValue = score
			}

		}
		if maxValue <= 0 {
			break
		} else {
			resultScores = append(resultScores, maxValue)
			resultPageIDs = append(resultPageIDs, pageIDs[maxIndex])
		}
		scores = append(scores[:maxIndex], scores[maxIndex+1:]...)
		pageIDs = append(pageIDs[:maxIndex], pageIDs[maxIndex+1:]...)
	}
	elapse := time.Since(Start)
	fmt.Printf("Time used in search=%v\n", elapse)
	db.Close()
	// elapse := time.Since(Start)
	// fmt.Printf("Time used in search=%v", elapse)
	return
}

func GetDLen(page SPage) (result float64) {
	result = 0
	for _, word := range page.myWord {
		result += word.TFIDF * word.TFIDF
	}
	return
}
func GetMaxTFByWordID(allPages []*spider.Page, wordID int) (MaxTF int) {
	MaxTF = 0
	for _, page := range allPages {
		for _, word := range page.Words() {
			if wordID == word.WordID {
				if word.TF() > MaxTF {
					MaxTF = word.TF()
					//MaxTF = 1
					//fmt.Printf("pageid:%v   wordID:%v    TF:%v\n", page.PageID, word.WordID, word.TF())
				}
			}
		}
	}

	return
}

func GetDFPhased(searchPages []*spider.Page, phaseterms []string) (result float64) {
	result = 0.0
	for _, page := range searchPages {
		if GetTFPhased(page, phaseterms) != 0.0 {
			result += 1
			//fmt.Printf("Page Url: %v", page.URL)
		}
	}
	return
}

func GetMaxTF(page *spider.Page) (maxTF float64) {
	maxTF = 0.0
	for _, word := range page.Words() {
		if float64(word.TF()) > maxTF {
			maxTF = float64(word.TF())
		}
	}
	return
}

func GetMaxTFPhased(allPage []*spider.Page, phaseterms []string) (MaxTF float64) {
	MaxTF = 0.0
	for _, page := range allPage {
		TF := GetTFPhased(page, phaseterms)
		if MaxTF < TF {
			MaxTF = TF
		}
	}
	return

}

func GetTFPhased(page *spider.Page, phaseterms []string) float64 {
	words := page.Words()
	var requiredWords []*spider.Word
	for _, phase := range phaseterms {
		for _, word := range words {
			if word.Word == phase {
				requiredWords = append(requiredWords, word)
			}
		}
	}
	// Doc doesn't contain all terms
	if len(phaseterms) != len(requiredWords) {
		return 0.0
	}
	resultPos := requiredWords[0].Positions()
	for index, requiredWord := range requiredWords {
		if index != 0 {
			tempPos := requiredWord.Positions()
			var tempResult []int
			for _, pos := range resultPos {
				//If the following pos don't have pos+1, remove the item
				if containInt(tempPos, pos+1) == true {
					tempResult = append(tempResult, pos)
				}
			}
			resultPos = tempResult
		}
	}
	// if len(resultPos) != 0 {
	// 	fmt.Printf("Terms:%v	%v\n", requiredWords[0].Positions(), requiredWords[1].Positions())
	// }

	return float64(len(resultPos))
}
func GetTFIDFPhased(page *spider.Page, allStoredPages []*spider.Page, invertedTable map[int][]int64, phaseterms []string) float64 {
	words := page.Words()
	TF := GetTFPhased(page, phaseterms)
	var docLen float64 = 0.0
	for _, word := range words {
		docLen += float64(word.TF())
	}
	N := float64(len(allStoredPages))
	df := GetDFPhased(allStoredPages, phaseterms)

	//k1 := 2.0
	//b := 0.75
	//firstTerm := (math.Log((N - df + 0.5) / (df + 0.5)))
	//secondTerm := ((k1 + 1) * TF / ((k1*(1-b) + b*docLen/AveDocLen) + TF))
	MaxTF := GetMaxTF(page)
	if MaxTF > 0 {
		//fmt.Print("TF >0")
	}
	if MaxTF <= 0 || df <= 0 {
		return 0
	}
	firstTerm := TF / MaxTF
	secondTerm := math.Log2(N / df)

	return firstTerm * secondTerm
}

func containInt(list []int, item int) bool {
	for _, temp := range list {
		if temp == item {
			return true
		}
	}
	return false
}
