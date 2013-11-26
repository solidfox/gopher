package main

import (
	"fmt"
	"gopher/spider"
	"time"
)

func main() {
	start := time.Now()

	db := spider.NewDBM("DBM.db")
	defer db.Close()
	pages := spider.Get300Pages()
	db.StorePages2(pages)
	elapse := time.Since(start)
	fmt.Printf("Time:%v", elapse)

}
