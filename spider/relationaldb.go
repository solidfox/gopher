package spider

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type RelationalDB struct {
	db        *sql.DB
	pageCache *Page
}

func NewRelationalDB(dbpath string) *RelationalDB {
	_, fileLoadErr := os.Stat(dbpath)
	dbDidNotExist := os.IsNotExist(fileLoadErr)

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	if dbDidNotExist {
		db.Exec(CreatePageInfoTableStmt)
		db.Exec(CreateWordsTableStmt)
		db.Exec(CreateLinksTableStmt)
	}

	return &RelationalDB{
		db: db,
	}
}

//Removes all data from DB
func (rdb *RelationalDB) Clear() {
	rdb.db.Exec("DELETE from pageInfo WHERE pageID >= 0")
	rdb.db.Exec("DELETE from words WHERE wordID >= 0")
}

func (rdb *RelationalDB) Close() {
	rdb.db.Close()
}
