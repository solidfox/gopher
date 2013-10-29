package main

import (
	"fmt"
	"gopher/spider"
	"time"
)

func main() {
	start := time.Now()
	db := spider.NewDBM("DBM.db")
	pages := spider.Get30Pages()
	db.StorePages2(pages)
	//db.DisplayInvertedTable()
	//db.GetPages2()

	//i := 0
	//fmt.Printf("-----------------------------------------------\n")
	pages2 := db.GetPages2()

	spider.PrintEntireIndex(pages2)
	// for _, p := range pages2 {
	// 	i++
	// 	fmt.Printf("PageID: %v\n", p.PageID)
	// 	fmt.Printf("PageSize: %v\n", p.Size)
	// 	fmt.Printf("PageTitle: %v\n", p.Title)
	// 	fmt.Printf("PageURL: %v\n", p.URL)
	// 	fmt.Printf("PageModified: %v\n", p.Modified)
	// 	fmt.Printf("PageWord: \n")
	// 	for _, word := range p.Words() {
	// 		fmt.Printf("%v", word.Word)
	// 		for _, pos := range word.Positions() {
	// 			fmt.Printf(" %v", pos)
	// 		}
	// 	}
	// 	fmt.Printf("***********************************************\n")
	// }

	elapsed := time.Since(start)
	fmt.Printf("Time spent: %v\n", elapsed)
	//fmt.Printf("\nnumberofpage:%v", i)
	db.Close()
}
