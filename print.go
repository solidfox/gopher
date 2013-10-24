package main

import (
	"fmt"
	"gopher/spider"
	//"runtime"
)

func main() {
	/*	pageChan := spider.GetPages()
		i := 0
		runtime.GOMAXPROCS(runtime.NumCPU())

		for {
			fmt.Println("waiting for pages")
			fmt.Println((<-pageChan).Words())
			i++
			fmt.Println(i)
			fmt.Print("Goroutines: ")
			fmt.Println(runtime.NumGoroutine())
		}*/

	table := indexHandler
	table.InitialIndexAndMaps()
	SetTablesFromDB(table)
	table.PrintEntrieIndex()
}
