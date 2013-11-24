package ranker

import (
	"fmt"
	"gopher/spider"
	"math"
	"strconv"
	"strings"
	//"time"
)

const (
	TFIDF = iota //=0
	Cos          //1
	//add here
)

var AveDocLen float64 = 0.0
var TotalDocLen int = 0

func contain(list []int64, item int64) bool {
	for _, temp := range list {
		if temp == item {
			return true
		}
	}
	return false
}

func containInt(list []int, item int) bool {
	for _, temp := range list {
		if temp == item {
			return true
		}
	}
	return false
}

func computeAveDocLen(pages []*spider.Page) {
	for _, page := range pages {
		words := page.Words()
		for _, word := range words {
			TotalDocLen += word.TF()
		}
	}
	AveDocLen = float64(TotalDocLen) / float64(len(pages))
}

func SearchingResult(query *spider.Page, option int) (resultDocIDs []int64) {
	db := spider.NewDBM("DBM.db")
	//get pages and inverted index
	invertedindex := db.GetInvertedIndex()
	allStoredPages := db.GetPages2()
	computeAveDocLen(allStoredPages)
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

	//reform query
	words := query.Words()
	var phases []string
	//words contains phase,
	words = words[:10]
	//cosWord := words
	words[0].Word = "kwun"
	words[1].Word = "chiu"
	words[0].Word = words[0].Word + " " + words[1].Word
	searchWords := words
	relationalDb := spider.NewRelationalDB("sqlite.db")
	// Add phase, tokenize it for searching(with word and wordID)
	for index, word := range searchWords {
		if strings.Contains(word.Word, " ") {
			phase := word.Word
			phases = append(phases, phase)
			//remove phase
			searchWords = append(searchWords[:index], searchWords[index+1:]...)
			//remove phase
			words = append(words[:index], words[index+1:]...)
			singleWords := strings.Fields(phase)
			for _, singleWord := range singleWords {
				newWord := spider.NewWord(singleWord)
				//newWord.WordID = int(relationalDb.WordIDOf(singleWord))
				searchWords = append(searchWords, newWord)
			}
		}
	}
	//Add wordID to all word
	var docIDs []int64
	var searchPage []*spider.Page
	for _, searchword := range searchWords {
		searchword.WordID = int(relationalDb.WordIDOf(searchword.Word))
		for _, id := range invertedTable[searchword.WordID] {
			if contain(docIDs, id) == false {
				docIDs = append(docIDs, id)
			}
		}
	}

	for _, docID := range docIDs {
		for _, page := range allStoredPages {
			if page.PageID == docID {
				searchPage = append(searchPage, page)
			}
		}
	}

	for _, word := range words {
		word.WordID = int(relationalDb.WordIDOf(word.Word))
	}

	fmt.Printf("\nPhase:%v \n", phases)

	docResults := make([]float64, len(searchPage))
	fmt.Printf("Length of docID: %v\n", len(docIDs))
	fmt.Printf("Length of searchPage: %v\n", len(searchPage))
	switch option {
	default:
		fmt.Errorf("Search with invalid option")
	case TFIDF:
		for index, page := range searchPage {
			docResults[index] = 0
			//non-phase part
			for _, word := range words {
				wordID := word.WordID
				docResults[index] += GetTFIDF(page, allStoredPages, invertedTable, wordID)
			}
			//phase part
			for _, phase := range phases {
				//var phaseTerm []spider.Word
				phaseTerms := strings.Fields(phase)
				// for _, phaseTerm := range phaseTerms{
				// 	phaseTerm
				// }
				temp := GetTFIDFPhased(page, allStoredPages, invertedTable, phaseTerms)
				docResults[index] += temp
				if temp != 0.0 {
					fmt.Print("Phase Search Success\n")
					fmt.Printf("temp: %v", temp)
				}
			}

			//fmt.Printf("%v", docResults)
		}
	case Cos:
		for index, page := range searchPage {
			docResults[index] = CosSim(page, allStoredPages, invertedTable, words, phases)
		}
	}

	relationalDb.Close()
	db.Close()
	return
}

func CosSim(page *spider.Page, allStoredPages []*spider.Page, invertedTable map[int][]int64, query []*spider.Word, phases []string) float64 {
	var dqSum float64 = 0
	var dLen float64 = 0
	var qLen float64 = 0
	qLen = float64(len(query)) + float64(len(phases))

	//Word part
	for _, word := range page.Words() {
		//compute Doc length
		dLen += (float64(word.TF())) * (float64(word.TF()))
		for _, queryword := range query {
			if word == queryword {
				dqSum += float64(word.TF()) * 1.0
			}
		}
	}
	//phase part
	for _, phase := range phases {
		//var phaseTerm []spider.Word
		phaseTerms := strings.Fields(phase)
		dqSum += GetTFPhased(page, phaseTerms) * 1.0
	}

	dLen = math.Sqrt(dLen)
	qLen = math.Sqrt(qLen)
	return dqSum / (dLen * qLen)
}

func GetTFIDF(page *spider.Page, allStoredPages []*spider.Page, invertedTable map[int][]int64, wordID int) float64 {
	words := page.Words()
	var TF float64 = 0.0
	var docLen float64 = 0.0
	for _, word := range words {
		if word.WordID == wordID {
			TF = float64(word.TF())
		}
		docLen += float64(word.TF())
	}
	N := float64(len(allStoredPages))
	df := float64(len(invertedTable[wordID]))
	k1 := 2.0
	b := 0.75
	firstTerm := (math.Log((N - df + 0.5) / (df + 0.5)))
	secondTerm := ((k1 + 1) * TF / ((k1*(1-b) + b*docLen/AveDocLen) + TF))
	return firstTerm * secondTerm
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
	k1 := 2.0
	b := 0.75
	firstTerm := (math.Log((N - df + 0.5) / (df + 0.5)))
	secondTerm := ((k1 + 1) * TF / ((k1*(1-b) + b*docLen/AveDocLen) + TF))
	return firstTerm * secondTerm
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
	if len(resultPos) != 0 {
		fmt.Printf("Terms:%v	%v\n", requiredWords[0].Positions(), requiredWords[1].Positions())
	}

	return float64(len(resultPos))
}
