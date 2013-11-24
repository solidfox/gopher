package spider

import (
	"gopher/spider"
	"os"
	"testing"
)

const (
	testdbfilename = "testsqlite.db"
)

func TestInsertPages(t *testing.T) {
	pages := []*spider.Page{spider.NewPage(), spider.NewPage()}
	pages[0].URL = "http://www.apple.com/"
	pages[1].URL = "http://www.google.com/"
	rdb := spider.NewRelationalDB(testdbfilename)
	rdb.InsertPagesAndSetIDs(pages)
	for i, page := range pages {
		if page.PageID != int64(i+1) {
			t.Log(page.PageID)
			t.Fail()
		}
	}
	os.Remove(testdbfilename)
}
