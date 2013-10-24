package spider

import (
	"database/sql"
	//"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

/*
CREATE TABLE 'pageInfo' (
        'pageID' INTEGER PRIMARY KEY AUTOINCREMENT,
        'size' INTEGER,
        'url' TEXT,
        'modifiedDate' DATETIME,
        'title' TEXT,
        'childLinks' TEXT
);
*/

const (
	DateFormat = "2006-01-02 15:04:05"
)

type PageInfoDB struct {
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

func NewPageInfoDB(dbpath string) (pageDB *PageInfoDB) {
	db, err := sql.Open("sqlite3", dbpath)
	db.Exec("DELETE from pageInfo WHERE pageID >= 0")
	if err != nil {
		panic(err)
	}
	pageDB = &PageInfoDB{
		db,
	}
	return pageDB
}

func (pageInfoDB *PageInfoDB) InsertPageInfo(pInfo *PageInfo) {
	pageInfoDB.db.Exec(
		"INSERT INTO 'pageInfo' "+
			"VALUES (?, ?, ?, datetime(?), ?, ?)",
		pInfo.id, pInfo.size, pInfo.url, pInfo.date, pInfo.title, pInfo.childLinks)
}

func (pageInfoDB *PageInfoDB) GetPageInfoByPageID(pageID int64) (pInfo *PageInfo) {
	// var datestring string
	pInfo = &PageInfo{
		0,
		0,
		"",
		time.Now(),
		"",
		"",
	}
	rows, err := pageInfoDB.db.Query(
		"SELECT * FROM pageInfo WHERE pageID = ?", pageID)
	if err == nil && rows.Next() {
		rows.Scan(&pInfo.id, &pInfo.size, &pInfo.url, &pInfo.date, &pInfo.title, &pInfo.childLinks)
		// fmt.Println(datestring)
		// pInfo.date, _ = time.Parse(DateFormat, datestring)
		return pInfo
	} else {
		return nil
	}
}
