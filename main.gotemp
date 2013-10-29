package main

import (
	//"fmt"
	"gopher/spider"
	"runtime"
)

func main() {
	pages := spider.Get30Pages()
	//i := 0
	runtime.GOMAXPROCS(runtime.NumCPU())

	/*	for {
			fmt.Println("waiting for pages")
			fmt.Println((<-pageChan).Words())
			i++
			fmt.Println(i)
			fmt.Print("Goroutines: ")
			fmt.Println(runtime.NumGoroutine())
		}
	*/
	table := spider.NewIndexHandler()
	table.StorePages(pages)

	spider.StorePagesToInvTable(table)
	spider.StorePagesToForwardTable(table)
	spider.StorePagesToPageInfo(table) // set wordMap too
	//spider.SetWord(table)
	//spider.DebugFowardtable()
	//spider.DebugInvertedtable()
}
