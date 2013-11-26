package spider

import (
	"fmt"
	"github.com/cznic/exp/dbm"
	//"gopher/ranker"
	"encoding/json"
	// "math"
	"os"
	// "strconv"
	// "strings"
	// "time"
)

const DBMname = "DBM.db"

var TotalDocLen float64 = 0
var o = &dbm.Options{}

type DBM struct {
	db  *dbm.DB
	rdb *RelationalDB
}

func NewDBM(name string) (d *DBM) {
	d = &DBM{
		db:  dbConnect(name),
		rdb: NewRelationalDB("r" + name),
	}
	return d
}

func (d *DBM) RDB() *RelationalDB {
	return d.rdb
}

func checkDbExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		//fmt.Printf("no such file or directory: %s", filename)
		return false
	} else {
		return true
	}
}

func dbConnect(name string) *dbm.DB {
	if checkDbExist(name) != true {
		db, err := dbm.Create(name, o)
		checkErr("Error: dbm can't create\n", err)
		return db
	} else {
		db, err := dbm.Open(name, o)
		checkErr("Error: dbm can't open\n", err)
		return db
	}
	return nil
}

func (d *DBM) Close() {
	mydb := d.db
	d.rdb.Close()
	mydb.Close()
}

func (d *DBM) StorePages(pages []*Page) {
	//Relational DB
	relationalDb := d.rdb
	relationalDb.InsertPagesAndSetIDs(pages)
	mydb := d.db
	//Key value DBinit
	fowardtable, err := mydb.Array("fowardtable")
	if err != nil {
		fmt.Printf("Error:forwardtable disconnected")
		panic(err)
	}
	invertedtable, err := mydb.Array("invertedtable")
	if err != nil {
		fmt.Printf("Error:invertedtable disconnected")
		panic(err)
	}

	invertedindex := make(map[int]*jsonDocList)

	for _, page := range pages {
		//forwardtable
		wordList := jsonWordList{}
		key := page.PageID
		words := page.Words()
		wordList.Words = make([]jsonWord, len(words))
		relationalDb.InsertWordsAndSetIDs(words)
		for i, word := range words {
			wordList.Words[i].Id = word.WordID
			wordList.Words[i].Pos = word.Positions()
		}
		jsonWords, err := json.Marshal(wordList)
		checkErr("Marshal failed:", err)
		fowardtable.Set(jsonWords, key)
		// input, _ := fowardtable.Get(key)
		// fmt.Printf("Input: %v\n", input)

		//invertedtable
		for _, word := range words {
			doc := jsonDoc{page.PageID, word.Positions()}
			var docList *jsonDocList
			docList, exists := invertedindex[word.WordID]
			if !exists {
				docList = &jsonDocList{}
				invertedindex[word.WordID] = docList
			}
			docList.addDoc(doc)
		}
	}
	for wordID, docs := range invertedindex {
		jsonDocs, err := json.Marshal(docs)
		checkErr("Marshal of json docs:", err)
		invertedtable.Set(jsonDocs, wordID)
	}
}

func (d *DBM) getDocsWithWordID(wordId int) (docs *jsonDocList) {
	mydb := d.db
	invertedtable, err := mydb.Array("invertedtable")
	if err != nil {
		fmt.Printf("Error:inverted table disconnected")
		//panic(err)
	}

	jsonDocs, _ := invertedtable.Get(wordId)
	if jsonDocs != nil {
		bytes := jsonDocs.([]byte)
		json.Unmarshal(bytes, &docs)
	}

	return docs
}

func (d *DBM) GetPagesForQuery(query *Page) (pageSlice []*Page, scoreSlice []float32) {

	terms := query.Words()
	d.rdb.InsertWordsAndSetIDs(terms)
	dfs := make([]int, len(terms))
	pageMap := make(map[int64]*Page)
	pageSlice = make([]*Page, 0, len(terms))

	for i, term := range terms {
		docList := d.getDocsWithWordID(term.WordID)
		dfs[i] = len(docList.Docs)
		for _, doc := range docList.Docs {
			// The word is a UNIQUE object for each page because the positions varies by page.
			word := NewWord(term.Word)
			word.positions = doc.Pos
			var page *Page
			page, exists := pageMap[doc.Id]
			if !exists {
				page = NewPage()
				pageSlice = append(pageSlice, page)
			}
			page.PageID = doc.Id
			page.AddWord(word)
		}
	}

	d.rdb.CompleteThePageInfoOf(pageSlice)

	n := d.rdb.PageCount()
	scoreSlice = make([]float32, len(pageSlice))
	for i, page := range pageSlice {
		scoreSlice[i] = cosineSimilarity(query, page, dfs, n)
	}

	return
}

func cosineSimilarity(query *Page, page *Page, df []int, n int64) float32 {
	tfxidfs := make([]float64, len(query.Words()))
	for i, word := range query.Words() {
		tfxidfs[i] = page.TFxIDF(word.Word, df[i], n)
	}
	dotProduct := float64(0)
	for _, tfxidf := range tfxidfs {
		dotProduct += tfxidf
	}
	cosineSimilarity := float32(dotProduct / (query.VectorLength() * page.VectorLength()))
	return cosineSimilarity
}

// func (d *DBM) GetTF(wordId int, docId int64) (tf int) {
// 	//TODO
// }

// func (d *DBM) GetTfidf(wordId int, docId int64) (weight float64) {
// 	N := float64(d.GetDocumentNumber())
// 	df := float64(d.Getdf(wordId))
// 	k1 := 2.0
// 	b := 0.75
// 	firstTerm := (math.Log((N - df + 0.5) / (df + 0.5)))
// 	secondTerm := ((k1 + 1) * float64(d.GetTf(wordId, docId))) / ((k1*(1-b) + b*float64(d.docLength(docId))/d.aveDocLength()) + float64(d.GetTf(wordId, docId)))
// 	return firstTerm * secondTerm
// }

// // This method is super slow
// // So, i make it private to prevent misuse
// func (d *DBM) cosSim(docId int64, query []string) float64 {

// 	var dqSum float64 = 0
// 	var dLen float64 = 0
// 	var qLen float64 = 0
// 	qLen = float64(len(query))
// 	relationalDb := NewRelationalDB("sqlite.db")

// 	mydb := d.db
// 	fowardtable, err := mydb.Array("fowardtable")
// 	if err != nil {
// 		fmt.Printf("Error:forward table disconnected")
// 		//panic(err)
// 	}

// 	temp, _ := fowardtable.Get(docId)
// 	longstr := temp.(string)
// 	statements := strings.Split(longstr, ";")
// 	start := time.Now()
// 	fmt.Printf("Number of statements: %v\n", len(statements))
// 	fmt.Printf("Number of query: %v\n", len(query))
// 	for _, statement := range statements {
// 		tokens := strings.Fields(statement)

// 		if len(tokens) != 0 {
// 			temp, _ := strconv.ParseInt(tokens[1], 10, 64)
// 			wordId, _ := strconv.ParseInt(tokens[0], 10, 64)
// 			//fmt.Printf("%v\n", temp)
// 			dLen += (float64(temp)) * (float64(temp))
// 			for _, word := range query {
// 				ID := relationalDb.WordIDOf(word)
// 				if ID == wordId {
// 					dqSum += float64(temp) * 1.0
// 				}
// 			}
// 		}
// 	}
// 	elapsed := time.Since(start)
// 	fmt.Printf("Time spent in func CosSim For loop: %v\n", elapsed)
// 	dLen = math.Sqrt(dLen)
// 	qLen = math.Sqrt(qLen)

// 	relationalDb.Close()
// 	//mydb.Close()
// 	//fmt.Printf("%v", ID)

// 	return dqSum / (dLen * qLen)
// }
