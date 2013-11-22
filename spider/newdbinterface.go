package spider

import (
	"fmt"
	"github.com/cznic/exp/dbm"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	//"time"
)

var TotalDocLen float64 = 0
var o = &dbm.Options{}

type DBM struct {
	db *dbm.DB
}

func NewDBM(name string) (d *DBM) {
	mydb := dbConnect(name)

	d = &DBM{
		mydb,
	}
	return d
}

const DBMname = "DBM.db"

func (d *DBM) storePages(pages []*Page) {
	//Relational DB
	relationalDb := NewRelationalDB("sqlite.db")
	relationalDb.InsertPagesAndSetIDs(pages)
	for _, page := range pages {
		fmt.Printf("%v\n", page.PageID)
	}
	mydb := d.db
	//init
	//forwardtable
	fowardtable, err := mydb.Array("fowardtable")
	if err != nil {
		fmt.Printf("Error:forwardtable disconnected")
		panic(err)
	}
	//invertedtable
	// invertedtable, err := mydb.Array("invertedtable")
	// if err != nil {
	// 	fmt.Errorf("Error:invertedtable disconnected")
	// }

	//invertedtable map
	// invertedindex := make(map[int]string)
	//storing for each page
	for _, page := range pages {
		//forwardtable
		key := page.PageID
		words := page.Words()
		relationalDb.InsertWordsAndSetIDs(words)
		var value string
		for _, word := range words {
			value += strconv.Itoa(word.WordID) + " " + strconv.Itoa(word.TF())
			positions := word.Positions()
			for _, pos := range positions {
				value += " " + strconv.Itoa(pos)
			}
			value = value + ";"
		}
		//put forwardtable in db
		fowardtable.Set(value, key)

		//invertedtable
		// for _, word := range words {
		// 	if invertedindex[word.WordID] == "" {
		// 		invertedindex[word.WordID] += strconv.FormatInt(page.PageID, 10)
		// 	} else {
		// 		invertedindex[word.WordID] += "" + strconv.FormatInt(page.PageID, 10)
		// 	}
		// }

	}
	//put invertedtable in db
	// for wordId, resultstr := range invertedindex {
	// 	invertedtable.Set(resultstr, wordId)
	// }
	relationalDb.Close()
	mydb.Close()
}

func (d *DBM) StorePages2(pages []*Page) {
	//fmt.Printf("numofpagestored%v\n", len(pages))
	//Relational DB
	relationalDb := NewRelationalDB("sqlite.db")
	relationalDb.InsertPagesAndSetIDs(pages)
	// for _, page := range pages {
	// 	fmt.Printf("%v\n", page.PageID)
	// }
	mydb := d.db
	//init
	//forwardtable
	fowardtable, err := mydb.Array("fowardtable")
	if err != nil {
		fmt.Printf("Error:forwardtable disconnected")
		panic(err)
	}

	//invertedtable
	invertedtable, err := mydb.Array("invertedtable")
	if err != nil {
		fmt.Printf("Error:invertedtable disconnected")
		panic(err)
	}

	//invertedtable map
	invertedindex := make(map[int]string)

	//fowardtable.Set("value", "key")
	for _, page := range pages {
		//forwardtable
		key := page.PageID
		words := page.Words()
		relationalDb.InsertWordsAndSetIDs(words)
		var value string = ""
		//fmt.Printf("numofwordinonepage:%v\n", len(words))
		for _, word := range words {
			value += strconv.Itoa(word.WordID) + " " + strconv.Itoa(word.TF())
			positions := word.Positions()
			for _, pos := range positions {
				value += " " + strconv.Itoa(pos)
			}
			value = value + ";"
			//fmt.Printf("%v		", value)
		}
		fowardtable.Set(value, key)
		// input, _ := fowardtable.Get(key)
		// fmt.Printf("Input: %v\n", input)

		//invertedtable
		for _, word := range words {
			if invertedindex[word.WordID] == "" {
				invertedindex[word.WordID] += strconv.FormatInt(page.PageID, 10)
			} else {
				invertedindex[word.WordID] += ";" + strconv.FormatInt(page.PageID, 10)
			}
		}

	}
	for key, value := range invertedindex {
		invertedtable.Set(value, key)
	}

	relationalDb.Close()
	//mydb.Close()
}

func (d *DBM) getPages() (pages []*Page) {
	relationalDb := NewRelationalDB("sqlite.db")
	mydb := d.db
	//forward table
	fowardtable, err := mydb.Array("fowardtable")
	if err != nil {
		//fmt.Printf("Error:forward table disconnected")
		//panic(err)
	}
	enum, err := fowardtable.Enumerator(true)
	if err != io.EOF {
		//fmt.Errorf("Error:fowardtable enum no exist")
		//panic(err)
	}
	key, value, err := enum.Next()
	if err != io.EOF {
		//fmt.Errorf("Error: fowardtable enum contain nothing")
		//panic(err)
	}

	for ; err != io.EOF; key, value, err = enum.Next() {
		page := NewPage()
		pid := key[0].(int64)
		page.PageID = pid
		words := make([]*Word, 0)
		valueindb := value[0].(string)
		str := strings.Split(valueindb, ";")
		for _, s := range str {
			token := strings.Fields(s)
			var word *Word
			for pos, val := range token {
				if pos == 0 {
					wordId, _ := strconv.Atoi(val)
					word.WordID = wordId
					name := relationalDb.WordOf(wordId)
					word.Word = name
				} else if pos == 1 {
					//TF :=strconv.Atoi(val)
					//word.
				} else {
					wordPos, _ := strconv.Atoi(val)

					word.AddPositions([]int{wordPos})
				}
			}
			words = append(words, word)

		}

		page.AddWords(words)
		/*value
		token := strings.Split(value[0].(string), ";")
		for _, j := range token {

			//i:=0
			table.forwardIndex[int(pid)] = make(map[int][]int)
			var token2 []string = strings.Fields(j)
			positions := make([]int, 0, DefaultPositionsLength)
			wid := 0
			for i, temp := range token2 {
				//temp:=token[num]
				if i == 0 {
					wid, _ = strconv.Atoi(temp)
				} else {
					positionInt, _ := strconv.Atoi(temp)
					positions = append(positions, positionInt)
				}
			}
			table.forwardIndex[int(pid)][wid] = positions
		}*/
	}

	relationalDb.CompleteThePageInfoOf(pages)

	relationalDb.Close()
	//mydb.Close()
	return pages
}

func (d *DBM) Close() {
	mydb := d.db
	mydb.Close()
}

func (d *DBM) DisplayInvertedTable() {
	mydb := d.db
	invertedtable, err := mydb.Array("invertedtable")
	if err != nil {
		fmt.Printf("Error:invertedtable can't find")
		panic(err)
	}
	enum, err := invertedtable.Enumerator(true)
	if err != nil {
		panic(err)
	}
	key, value, err := enum.Next()
	if err != io.EOF {
		//fmt.Printf("Error:enum is empty")
		//panic(err)
	}
	for ; err != io.EOF; key, value, err = enum.Next() {
		integer := key[0].(int64)
		str := value[0].(string)
		fmt.Printf("%v:%v\n", integer, str)
	}
}

func (d *DBM) GetPages2() (pages []*Page) {
	relationalDb := NewRelationalDB("sqlite.db")
	mydb := d.db

	fowardtable, err := mydb.Array("fowardtable")
	if err != nil {
		fmt.Printf("Error:forward table disconnected")
		//panic(err)
	}

	enum, err := fowardtable.Enumerator(true)
	if err != nil {
		fmt.Printf("Error:fowardtable can't find enum")
		//panic(err)
	}
	//key, value, err := enum.Next()
	key, value, err := enum.Next()
	if err != io.EOF {
		//fmt.Printf("Error:enum is empty")
		//panic(err)
	}
	//fmt.Printf("%v    %v\n", key[0], value[0])
	for ; err != io.EOF; key, value, err = enum.Next() {
		page := NewPage()

		//relationalDb.CompleteThePageInfoOf(pages)
		//fmt.Printf("-----------------------------------------------\n")
		//fmt.Printf("%v    %v\n", key[0], value[0])
		//To keep it safe
		if key[0] != nil {
			page.PageID = key[0].(int64)
			//To keep it safe
			if value[0] != nil {
				//handle the long string
				longstr := value[0].(string)
				oneWordStatement := strings.Split(longstr, ";")
				var words []*Word
				for _, statement := range oneWordStatement {
					word := NewWord("")
					var wordId, wordTF int
					var pos []int
					items := strings.Fields(statement)
					for index, item := range items {
						if index == 0 {
							wordId, _ = strconv.Atoi(item)
						} else if index == 1 {

							wordTF, _ = strconv.Atoi(item)
							if wordTF == 10000 {

							}
						} else {
							temp, _ := strconv.Atoi(item)
							pos = append(pos, temp)
						}
					}
					word.WordID = wordId
					word.Word = relationalDb.WordOf(wordId)
					word.positions = pos
					words = append(words, word)
				}
				page.AddWords(words)

			}

		}
		pages = append(pages, page)
	}
	relationalDb.CompleteThePageInfoOf(pages)
	relationalDb.Close()
	//mydb.Close()
	//pages = append(pages, NewPage())
	return pages
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
		if err != nil {
			fmt.Printf("Error: dbm can't create\n")
			panic(err)
		} else {
			return db
		}
	} else {
		db, err := dbm.Open(name, o)
		if err != nil {
			fmt.Printf("Error: dbm can't open\n")
			panic(err)
		} else {
			return db
		}
	}
	return nil
}

func (d *DBM) GetWordNumber() (score int64) {
	var wordN string
	relationalDb := NewRelationalDB("sqlite.db")
	row := relationalDb.db.QueryRow(
		"SELECT COUNT(word) FROM words")
	row.Scan(&wordN)
	score, _ = strconv.ParseInt(wordN, 10, 64)
	relationalDb.Close()
	return score
}
func (d *DBM) GetDocumentNumber() (score int) {
	var docN string
	relationalDb := NewRelationalDB("sqlite.db")
	row := relationalDb.db.QueryRow(
		"SELECT COUNT(PageID) FROM pageInfo")
	row.Scan(&docN)
	temp, _ := strconv.ParseInt(docN, 10, 64)
	score = int(temp)
	//fmt.Printf("\nWord with id=10: %v\n", relationalDb.WordOf(10))
	relationalDb.Close()
	return score
}

func GetWordNumberByID(wordID int) (score int64) {
	var wordN string
	relationalDb := NewRelationalDB("sqlite.db")
	row := relationalDb.db.QueryRow(
		"SELECT COUNT(word) FROM words WHERE wordID = ?", wordID)
	row.Scan(&wordN)
	score, _ = strconv.ParseInt(wordN, 10, 64)
	relationalDb.Close()
	return score
}

func (d *DBM) Getdf(wordId int) (score int) {
	score = 0
	mydb := d.db
	invertedtable, err := mydb.Array("invertedtable")
	if err != nil {
		fmt.Printf("Error:inverted table disconnected")
		//panic(err)
	}
	//fmt.Printf("%v\n", wordId)
	temp, _ := invertedtable.Get(wordId)
	if temp != nil {
		str := temp.(string)
		statements := strings.Split(str, ";")
		score = len(statements)
	}

	return score
}

func (d *DBM) GetDocIdByWordID(wordId int) (docIDs []int64) {
	mydb := d.db
	invertedtable, err := mydb.Array("invertedtable")
	if err != nil {
		fmt.Printf("Error:inverted table disconnected")
		//panic(err)
	}

	temp, _ := invertedtable.Get(wordId)
	if temp != nil {
		str := temp.(string)
		statements := strings.Split(str, ";")
		for _, statement := range statements {
			//fmt.Printf("%v\n", statement)
			pageId, _ := strconv.ParseInt(statement, 10, 64)
			docIDs = append(docIDs, pageId)
		}
	}

	// enum, err := invertedtable.Enumerator(true)
	// if err != nil {
	// 	fmt.Printf("Error:invertedtable can't find enum")
	// 	//panic(err)
	// }
	// //key, value, err := enum.Next()
	// key, value, err := enum.Next()
	// if err != io.EOF {
	// 	//fmt.Printf("Error:enum is empty")
	// 	//panic(err)
	// }

	// for ; err != io.EOF; key, value, err = enum.Next() {

	// 	if wordId == int(key[0].(int64)) {
	// 		str := value[0].(string)
	// 		Statements := strings.Split(str, ";")
	// 		//docIDs = 0
	// 		for _, statement := range Statements {
	// 			//fmt.Printf("%v\n", statement)
	// 			pageId, _ := strconv.ParseInt(statement, 10, 64)
	// 			docIDs = append(docIDs, pageId)
	// 		}
	// 	}
	// }
	return docIDs
}

func (d *DBM) GetAllDocId() (docIDs []int64) {
	mydb := d.db
	invertedtable, err := mydb.Array("invertedtable")
	if err != nil {
		fmt.Printf("Error:inverted table disconnected")
		//panic(err)
	}

	enum, err := invertedtable.Enumerator(true)
	if err != nil {
		fmt.Printf("Error:invertedtable can't find enum")
		//panic(err)
	}
	//key, value, err := enum.Next()
	_, value, err := enum.Next()
	if err != io.EOF {
		//fmt.Printf("Error:enum is empty")
		//panic(err)
	}

	for ; err != io.EOF; _, value, err = enum.Next() {

		//if wordId == int(key[0].(int64)) {
		str := value[0].(string)
		Statements := strings.Split(str, ";")
		//docIDs = 0
		for _, statement := range Statements {
			//fmt.Printf("%v\n", statement)
			pageId, _ := strconv.ParseInt(statement, 10, 64)
			docIDs = append(docIDs, pageId)
		}
		//}
	}
	return docIDs
}

func (d *DBM) GetTf(wordId int, docId int64) (tf int) {
	mydb := d.db
	fowardtable, err := mydb.Array("fowardtable")
	if err != nil {
		fmt.Printf("Error:forward table disconnected")
		//panic(err)
	}

	temp, _ := fowardtable.Get(docId)
	longstr := temp.(string)
	statements := strings.Split(longstr, ";")
	for _, statement := range statements {
		tokens := strings.Fields(statement)

		if len(tokens) != 0 {
			a, _ := strconv.ParseInt(tokens[0], 10, 64)
			wordId2 := int(a)

			if wordId == wordId2 {
				temp, _ := strconv.ParseInt(tokens[1], 10, 64)
				//fmt.Printf("%v\n", temp)
				tf = int(temp)
			}
		}
	}

	return tf
}

func (d *DBM) docLength(docId int64) (length int) {
	length = 0
	mydb := d.db
	fowardtable, err := mydb.Array("fowardtable")
	if err != nil {
		fmt.Printf("Error:forward table disconnected")
		//panic(err)
	}

	temp, _ := fowardtable.Get(docId)
	longstr := temp.(string)
	statements := strings.Split(longstr, ";")
	for _, statement := range statements {
		tokens := strings.Fields(statement)

		if len(tokens) != 0 {
			temp, _ := strconv.ParseInt(tokens[1], 10, 64)
			//fmt.Printf("%v\n", temp)
			length += int(temp)
		}
	}
	// for ; err != io.EOF; key, value, err = enum.Next() {
	// 	docId2 := key[0].(int64)
	// 	if docId == docId2 {
	// 		longstr := value[0].(string)
	// 		statements := strings.Split(longstr, ";")
	// 		for _, statement := range statements {
	// 			tokens := strings.Fields(statement)

	// 			if len(tokens) != 0 {
	// 				temp, _ := strconv.ParseInt(tokens[1], 10, 64)
	// 				//fmt.Printf("%v\n", temp)
	// 				length += int(temp)
	// 			}
	// 		}
	// 	}
	// }
	return length
}

func (d *DBM) aveDocLength() (length float64) {
	N := float64(d.GetDocumentNumber())
	if TotalDocLen == 0.0 {
		length = 0
		docIds := d.GetAllDocId()
		for _, docId := range docIds {
			length += float64(d.docLength(docId))
		}
		TotalDocLen = length
	} else {
		length = TotalDocLen
	}

	return (length / N)
}

func (d *DBM) GetTfidf(wordId int, docId int64) (weight float64) {
	N := float64(d.GetDocumentNumber())
	df := float64(d.Getdf(wordId))
	k1 := 2.0
	b := 0.75
	firstTerm := (math.Log((N - df + 0.5) / (df + 0.5)))
	secondTerm := ((k1 + 1) * float64(d.GetTf(wordId, docId))) / ((k1*(1-b) + b*float64(d.docLength(docId))/d.aveDocLength()) + float64(d.GetTf(wordId, docId)))
	return firstTerm * secondTerm
}

func (d *DBM) CosSim(docId int64, query []string) float64 {
	var dqSum float64 = 0
	var dLen float64 = 0
	var qLen float64 = 0
	qLen = float64(len(query))
	relationalDb := NewRelationalDB("sqlite.db")

	mydb := d.db
	fowardtable, err := mydb.Array("fowardtable")
	if err != nil {
		fmt.Printf("Error:forward table disconnected")
		//panic(err)
	}

	temp, _ := fowardtable.Get(docId)
	longstr := temp.(string)
	statements := strings.Split(longstr, ";")
	for _, statement := range statements {
		tokens := strings.Fields(statement)

		if len(tokens) != 0 {
			temp, _ := strconv.ParseInt(tokens[1], 10, 64)
			wordId, _ := strconv.ParseInt(tokens[0], 10, 64)
			//fmt.Printf("%v\n", temp)
			dLen += (float64(temp)) * (float64(temp))
			for _, word := range query {
				ID := relationalDb.WordIDOf(word)
				if ID == wordId {
					dqSum += float64(temp) * 1.0
				}
			}
		}
	}

	dLen = math.Sqrt(dLen)
	qLen = math.Sqrt(qLen)

	relationalDb.Close()
	//mydb.Close()
	//fmt.Printf("%v", ID)
	return dqSum / (dLen * qLen)
}
