package main

import (
	//"fmt"
	"gopher/spider"
)

func main() {
	db := spider.NewDBM("DBM.db")
	db.StorePages(spider.Get30Pages())
}
