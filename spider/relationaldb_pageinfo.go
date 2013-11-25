package spider

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	CreatePageInfoTableStmt = `
	CREATE TABLE pageInfo (
		pageID 			INTEGER PRIMARY KEY AUTOINCREMENT,
		size 			INTEGER,
		url 			TEXT UNIQUE,
		modifiedDate 	DATETIME,
		title 			TEXT
	)`
	CreateLinksTableStmt = `
	CREATE TABLE links (
		parent		INTEGER REFERENCES pageInfo(pageID),
		child		INTEGER REFERENCES pageInfo(pageID),
		UNIQUE(parent, child))`

	InsertPageStmt = `
	INSERT OR IGNORE INTO 'pageInfo' 
	VALUES (NULL, ?, ?, ?, ?, ?)`
	UpdatePageStmt = `
	UPDATE 'pageInfo' SET 
		size = ?
		modifiedDate = ?
		title = ?
		WHERE url = ?`
	GetPageIdStmt        = `SELECT pageID FROM pageInfo WHERE url = ?`
	GetPagesByPageIDStmt = `SELECT * FROM pageInfo WHERE pageID = ?`
	GetPagesByURLStmt    = `SELECT * FROM pageInfo WHERE url = ?`

	InsertLinkStmt = `
	INSERT OR IGNORE INTO links 
	VALUES (?, ?)`
	GetParentsStmt = `
	SELECT title, url 
	FROM pageInfo INNER JOIN links ON pageInfo.pageID = links.parent
	WHERE links.child = ?`
	GetChildrenStmt = `
	SELECT title, url 
	FROM pageInfo INNER JOIN links ON pageInfo.pageID = links.child
	WHERE links.parent = ?`
)

// func (rdb *RelationalDB) InsertLink(parent string, child string) {
// 	tx, _ := rdb.db.Begin()
// 	defer tx.Commit()
// 	insertLinkStmt := tx.Prepare(InsertLinkStatement)

// 	insertLink(stmt, parent, child)
// }

func (rdb *RelationalDB) insertLink(
	insertLinkStmt *sql.Stmt,
	getPageIDStmt *sql.Stmt,
	insertPageStmt *sql.Stmt,
	parent string,
	child string,
) {
	var parentID int64
	var childID int64
	if rdb.pageCache.URL == parent {
		parentID = rdb.pageCache.PageID
	} else {
		row := getPageIDStmt.QueryRow(parent)
		err := row.Scan(&parentID)
		if err == sql.ErrNoRows {
			insertPageStmt.Exec(nil, parent, nil, nil)
		}
	}
	insertPageStmt.Exec(nil, child, nil, nil)
	row := getPageIDStmt.QueryRow(child)
	err := row.Scan(&childID)
	if err != nil {
		log.Fatal(err)
	}
	insertLinkStmt.Exec(parentID, childID)
}

func (rdb *RelationalDB) InsertPagesAndSetIDs(pages []*Page) {
	tx, _ := rdb.db.Begin()
	defer tx.Commit()
	updatePage, _ := tx.Prepare(UpdatePageStmt)
	insertPage, _ := tx.Prepare(InsertPageStmt)
	getPageId, _ := tx.Prepare(GetPageIdStmt)
	insertLink, _ := tx.Prepare(InsertLinkStmt)

	for _, p := range pages {
		rdb.pageCache = p
		res, err := updatePage.Exec(p.Size, p.Modified, p.Title, p.URL)
		if err != nil {
			log.Fatal(err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		if rowsAffected == 0 {
			insertPage.Exec(p.Size, p.URL, p.Modified, p.Title)
		}
		row := getPageId.QueryRow(p.URL)
		row.Scan(&p.PageID)
		for _, link := range p.Links() {
			rdb.insertLink(insertLink, getPageId, insertPage, p.URL, link.URL)
			//TODO update capability: remove links no longer in the page
		}
	}
}

func (rdb *RelationalDB) PageByPageID(pageID int64) (p *Page) {
	p = NewPage()
	rdb.CompleteThePageInfoOf([]*Page{p})
	return p
}

// Uses the Page's PageID or URL to fill out all information about the page except the words it
// contains (since those are not held in the relational db).
func (rdb *RelationalDB) CompleteThePageInfoOf(pages []*Page) {
	tx, _ := rdb.db.Begin()
	defer tx.Commit()
	getPagesByPageID, _ := tx.Prepare(GetPagesByPageIDStmt)
	getPagesByURL, _ := tx.Prepare(GetPagesByURLStmt)
	for _, p := range pages {
		if p.PageID != -1 {
			row := getPagesByPageID.QueryRow(p.PageID)
			row.Scan(&p.PageID, &p.Size, &p.URL, &p.Modified, &p.Title)
		} else if len(p.URL) > 0 {
			row := getPagesByURL.QueryRow(p.URL)
			row.Scan(&p.PageID, &p.Size, &p.URL, &p.Modified, &p.Title)
		}
	}
}

func (rdb *RelationalDB) LoadParentsFor(pages []*Page) {
	parents := make([]*Link, 0, 10)
	tx, _ := rdb.db.Begin()
	defer tx.Commit()
	getParents, _ := tx.Prepare(GetParentsStmt)
	for _, p := range pages {
		if p.PageID < 0 {
			log.Println("RelationalDB.LoadParentsFor() requires passed pages to have PageIDs: ")
			log.Fatalln(p)
		}
		rows, err := getParents.Query(p.PageID)
		if err == nil {
			log.Fatalln(err)
		}
		for rows.Next() {
			link := &Link{}
			rows.Scan(&link.Title, &link.URL)
			parents = append(parents, link)
			p.Parents = parents
		}
	}
}

func (rdb *RelationalDB) LoadChildrenFor(pages []*Page) {
	children := make([]*Link, 0, 10)
	tx, _ := rdb.db.Begin()
	defer tx.Commit()
	getChildren, _ := tx.Prepare(GetChildrenStmt)
	for _, p := range pages {
		if p.PageID < 0 {
			log.Println("RelationalDB.LoadChildrenFor() requires passed pages to have PageIDs: ")
			log.Fatalln(p)
		}
		rows, err := getChildren.Query(p.PageID)
		if err == nil {
			log.Fatalln(err)
		}
		for rows.Next() {
			link := &Link{}
			rows.Scan(&link.Title, &link.URL)
			children = append(children, link)
			p.Parents = children
		}
	}
}
