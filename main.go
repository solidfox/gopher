package main

import (
	//"fmt"
	"gopher/spider"
)

func main() {

	db := spider.NewDBM("DBM.db")
	pages := spider.Get30Pages()
	db.StorePages(pages)

	db.Close()

}
