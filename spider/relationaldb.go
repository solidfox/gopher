package spider

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

/*
CREATE TABLE 'pageInfo' (
    'pageID' 		INTEGER PRIMARY KEY AUTOINCREMENT,
    'size' 			INTEGER,
    'url' 			TEXT UNIQUE,
    'modifiedDate' 	DATETIME,
    'title' 		TEXT,
    'childLinks' 	TEXT
);
CREATE TABLE 'words' (
	'word'		TEXT UNIQUE,
	'wordID'	INTEGER PRIMARY KEY AUTOINCREMENT
);
*/

const (
	DateFormat = "2006-01-02 15:04:05"
)

type RelationalDB struct {
	db *sql.DB
}

type PageInfo struct {
	id         int64
	size       int64
	url        string
	date       time.Time
	title      string
	childLinks string
}

func NewRelationalDB(dbpath string) (rdb *RelationalDB) {
	db, err := sql.Open("sqlite3", dbpath)
	// db.Exec("DELETE from pageInfo WHERE pageID >= 0")
	if err != nil {
		panic(err)
	}
	rdb = &RelationalDB{
		db,
	}
	return rdb
}

// //Removes all data from DB
// func (rdb *RelationalDB) Clear() {

// }

func (rdb *RelationalDB) InsertPageInfo(pInfo *PageInfo) {
	rdb.db.Exec(
		"INSERT OR IGNORE INTO 'pageInfo' "+
			"VALUES (?, ?, ?, datetime(?), ?, ?)",
		pInfo.id, pInfo.size, pInfo.url, pInfo.date, pInfo.title, pInfo.childLinks)
}

func (rdb *RelationalDB) InsertWord(word string) {
	rdb.db.Exec("INSERT OR IGNORE into words VALUES (?, DEFAULT)", word)
}

func (rdb *RelationalDB) WordIDOf(word string) (wordID int64) {
	rows, err := rdb.db.Query(
		"SELECT wordID FROM words WHERE word = ?", word)
	if err == nil && rows.Next() {
		rows.Scan(&wordID)
		rows.Close()
		return pInfo
	} else {
		return nil
	}
}

func (rdb *RelationalDB) WordOf(wordID int64) (word string) {
	rows, err := rdb.db.Query(
		"SELECT wordID FROM words WHERE wordID = ?", wordID)
	if err == nil && rows.Next() {
		rows.Scan(&wordID)
		rows.Close()
		return pInfo
	} else {
		return nil
	}
}

func (rdb *RelationalDB) Close() {
	rdb.db.Close()
}

func (rdb *RelationalDB) GetPageInfoByPageID(pageID int64) (pInfo *PageInfo) {
	// var datestring string
	pInfo = &PageInfo{
		0,
		0,
		"",
		time.Now(),
		"",
		"",
	}
	rows, err := rdb.db.Query(
		"SELECT * FROM pageInfo WHERE pageID = ?", pageID)
	if err == nil && rows.Next() {
		rows.Scan(&pInfo.id, &pInfo.size, &pInfo.url, &pInfo.date, &pInfo.title, &pInfo.childLinks)
		rows.Scan()
		return pInfo
	} else {
		return nil
	}
}
