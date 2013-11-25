package spider

import (
	"fmt"
	"github.com/cznic/exp/dbm"
	//"gopher/ranker"
	"encoding/json"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var TotalDocLen float64 = 0
var o = &dbm.Options{}

type DBM struct {
	db *dbm.DB
}

type jsonDocList struct {
	DF   int64
	docs []jsonDoc
}

type jsonDoc struct {
	id  int64
	tf  int
	pos []int
}

type jsonWordList struct {
	words []jsonWord
}

type jsonWord struct {
	id  int64
	pos []int
}

func NewDBM(name string) (d *DBM) {
	mydb := dbConnect(name)

	d = &DBM{
		mydb,
	}
	return d
}

const DBMname = "DBM.db"

func (d *DBM) StorePages(pages []*Page) {
	//Relational DB
	relationalDb := NewRelationalDB("sqlite.db")
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

	invertedindex := make(map[int]jsonDocList)
	docList := jsonDocList{}

	for _, page := range pages {
		//forwardtable
		wordList := jsonWordList{}
		key := page.PageID
		words := page.Words()
		wordList.words = make([]jsonWord, len(words))
		relationalDb.InsertWordsAndSetIDs(words)
		for i, word := range words {
			wordList.words[i].id = word.WordID
			wordList.words[i].pos = word.Positions()
		}
		json, err := json.Marshal(wordList)
		checkErr("Marshal failed:", err)
		fowardtable.Set(string(json), key)
		// input, _ := fowardtable.Get(key)
		// fmt.Printf("Input: %v\n", input)

		//invertedtable
		for _, word := range words {
			invertedindex[word.WordID].addDoc()
		}

	}
	for key, value := range invertedindex {
		invertedtable.Set(value, key)
	}

	relationalDb.Close()
	//mydb.Close()
}
