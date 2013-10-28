package spider

import (
	"fmt"
	"github.com/cznic/exp/dbm"
	"io"
	"os"
	"strconv"
	"strings"
)

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

func (d *DBM) StorePages(pages []*Page) {
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
	invertedtable, err := mydb.Array("invertedtable")
	if err != nil {
		fmt.Errorf("Error:invertedtable disconnected")
	}

	//invertedtable map
	invertedindex := make(map[int]string)
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
		for _, word := range words {
			if invertedindex[word.WordID] == "" {
				invertedindex[word.WordID] += strconv.FormatInt(page.PageID, 10)
			} else {
				invertedindex[word.WordID] += "" + strconv.FormatInt(page.PageID, 10)
			}
		}

	}
	//put invertedtable in db
	for wordId, resultstr := range invertedindex {
		invertedtable.Set(resultstr, wordId)
	}
	relationalDb.Close()
	mydb.Close()
}
func (d *DBM) GetPages() (pages []*Page) {
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
	mydb.Close()
	return pages
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
	key, value, err := enum.Next()
	if err != io.EOF {
		fmt.Printf("Error:enum is empty")
		//panic(err)
	}
	//fmt.Printf("%v    %v", key[0], value[0])
	for ; err != io.EOF; key, value, err = enum.Next() {
		relationalDb.CompleteThePageInfoOf(pages)
		fmt.Printf("-----------------------------------------------\n")
		fmt.Printf("%v    %v\n", key[0], value[0])
	}

	relationalDb.Close()
	mydb.Close()
	pages = append(pages, NewPage())
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
