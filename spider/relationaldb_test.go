package spider

import (
	"gopher/spider"
	"testing"
)

func TestInsertPages(t *testing.T) {
	// pages := []*spider.Page{spider.NewPage()}
	rdb := spider.NewRelationalDB("sqlite.db")
	rdb.InsertPagesAndSetIDs(spider.Get30Pages())
	rdb.Clear()
	t.Fail()
}
