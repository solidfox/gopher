package spider

import (
	"gopher/spider"
	"testing"
)

func TestInsertPages(t *testing.T) {
	rdb := spider.NewRelationalDB("sqlite.db")
	//rdb.InsertPagesAndSetIDs(spider.Get30Pages())
	rdb.Clear()
	t.Fail()
}
