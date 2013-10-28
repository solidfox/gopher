package spider

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	// "log"
	//"fmt"
	"os"
	"strings"
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

type RelationalDB struct {
	db *sql.DB
}

func NewRelationalDB(dbpath string) *RelationalDB {
	_, fileLoadErr := os.Stat(dbpath)
	dbDidNotExist := os.IsNotExist(fileLoadErr)

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	if dbDidNotExist {
		db.Exec(
			"CREATE TABLE pageInfo (" +
				"pageID 		INTEGER PRIMARY KEY AUTOINCREMENT," +
				"size 			INTEGER," +
				"url 			TEXT UNIQUE," +
				"modifiedDate 	DATETIME," +
				"title 			TEXT," +
				"childLinks 	TEXT" +
				")")
		db.Exec(
			"CREATE TABLE words (" +
				"word		TEXT UNIQUE," +
				"wordID	INTEGER PRIMARY KEY AUTOINCREMENT" +
				")")
	}

	return &RelationalDB{
		db,
	}
}

//Removes all data from DB
func (rdb *RelationalDB) Clear() {
	rdb.db.Exec("DELETE from pageInfo WHERE pageID >= 0")
	rdb.db.Exec("DELETE from words WHERE wordID >= 0")
}

func (rdb *RelationalDB) InsertPagesAndSetIDs(pages []*Page) {
	// _, err := rdb.db.Exec("INSERT INTO 'pageInfo' VALUES (NULL, 0, 'www.blue.com', 0, 'Bluebusters', 'mojgoj')")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	tx, _ := rdb.db.Begin()

	for _, p := range pages {
		links := make([]string, len(p.Links()))
		for i, link := range p.Links() {
			links[i] = link.URL
		}
		if p.PageID == -1 {
			tx.Exec(
				"INSERT OR IGNORE INTO 'pageInfo' "+
					"VALUES (NULL, ?, ?, ?, ?, ?)",
				p.Size, p.URL, p.Modified, p.Title, strings.Join(links, " "))
			row := tx.QueryRow("SELECT pageID FROM pageInfo WHERE url = ?", p.URL)
			row.Scan(&p.PageID)
		}
		// if i%3 == 0 {
		// 	fmt.Printf("We saved %v pages\n", i)
		// }

	}
	tx.Commit()
}

func (rdb *RelationalDB) InsertWordsAndSetIDs(words []*Word) {
	tx, _ := rdb.db.Begin()
	for _, w := range words {
		if w.WordID == -1 {
			tx.Exec("INSERT OR IGNORE INTO 'words' VALUES (?, NULL)", w.Word)
			row := tx.QueryRow("SELECT wordID FROM words WHERE word = ?", w.Word)
			row.Scan(&w.WordID)
		}
		// if i%1000 == 0 {
		// 	fmt.Printf("We saved %v words\n", i)
		// }
	}
	tx.Commit()
}

func (rdb *RelationalDB) Close() {
	rdb.db.Close()
}

func (rdb *RelationalDB) PageByPageID(pageID int64) (p *Page) {
	p = NewPage()
	rdb.CompleteThePageInfoOf([]*Page{p})
	return p
}

// Uses the Page's PageID to fill out all information about the page except the words it
// contains (since those are not held in the relational db).
func (rdb *RelationalDB) CompleteThePageInfoOf(pages []*Page) {
	var links string
	tx, _ := rdb.db.Begin()
	for _, p := range pages {
		if p.PageID != -1 {
			row := tx.QueryRow("SELECT * FROM pageInfo WHERE pageID = ?", p.PageID)
			row.Scan(&p.PageID, &p.Size, &p.URL, &p.Modified, &p.Title, &links)
			linkSlice := strings.Fields(links)
			for _, link := range linkSlice {
				p.AddLink(link, "")
			}
		}
	}
	tx.Commit()
}

// Fills out the Word's WordID provided that it's Word field is not the empty string.
func (rdb *RelationalDB) AddWordIDTo(words []*Word) {
	tx, _ := rdb.db.Begin()
	for _, w := range words {
		if w.Word != "" {
			row := tx.QueryRow("SELECT wordID FROM words WHERE word = ?", w.Word)
			row.Scan(&w.WordID)
		}
	}
	tx.Commit()
}

// Fills out the Word's Word field provided that it's WordID field is not -1.
func (rdb *RelationalDB) AddWordWordTo(words []*Word) {
	tx, _ := rdb.db.Begin()
	for _, w := range words {
		if w.WordID != -1 {
			row := tx.QueryRow("SELECT word FROM words WHERE wordID = ?", w.WordID)
			row.Scan(&w.WordID)
		}
	}
	tx.Commit()
}

func (rdb *RelationalDB) WordIDOf(word string) (wordID int64) {
	row := rdb.db.QueryRow(
		"SELECT wordID FROM words WHERE word = ?", word)
	row.Scan(&wordID)
	return wordID
}

func (rdb *RelationalDB) WordOf(wordID int) (word string) {
	row := rdb.db.QueryRow(
		"SELECT word FROM words WHERE wordID = ?", wordID)
	row.Scan(&word)
	return word
}
