package spider

import (
	_ "github.com/mattn/go-sqlite3"
)

const (
	CreateWordsTableStmt = `
	CREATE TABLE words (
		word		TEXT UNIQUE,
		wordID	INTEGER PRIMARY KEY AUTOINCREMENT
	)`
)

func (rdb *RelationalDB) InsertWordsAndSetIDs(words []*Word) {
	tx, _ := rdb.db.Begin()
	for _, w := range words {
		if w.WordID == -1 {
			tx.Exec("INSERT OR IGNORE INTO 'words' VALUES (?, NULL)", w.Word)
			row := tx.QueryRow("SELECT wordID FROM words WHERE word = ?", w.Word)
			row.Scan(&w.WordID)
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

func (rdb *RelationalDB) WordCount() (count int64) {
	row := rdb.db.QueryRow("SELECT COUNT(word) FROM words")
	row.Scan(&count)
	return count
}
