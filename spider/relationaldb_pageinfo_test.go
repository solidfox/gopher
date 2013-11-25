package spider

import (
	"os"
	"testing"
)

const (
	testdbfilename = "testsqlite.db"
)

var rdb *RelationalDB

func TestInsertPages(t *testing.T) {
	rdb = NewRelationalDB(testdbfilename)
	pages := []*Page{NewPage(), NewPage()}
	pages[0].URL = "http://www.apple.com/"
	pages[1].URL = "http://www.google.com/"
	rdb.InsertPagesAndSetIDs(pages)
	for i, page := range pages {
		if page.PageID != int64(i+1) {
			t.Log(page.PageID)
			t.Fail()
		}
	}
	rdb.InsertPagesAndSetIDs(pages)
	for i, page := range pages {
		if page.PageID != int64(i+1) {
			t.Log(page.PageID)
			t.Fail()
		}
	}
	os.Remove(testdbfilename)
}
