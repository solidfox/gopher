package main

import (
	"fmt"
	"gopher/spider"
)

func main() {

	db := spider.NewDBM("DBM.db")
	pages := db.GetPages2()
	//db.StorePages2(pages)
	words := pages[0].Words()
	for _, word := range words {
		fmt.Printf("%v\n", word)
	}

	db.Close()

}
