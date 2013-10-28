package main

import (
	"gopher/spider"
	//"runtime"
	"os"
	"strconv"
)

func PrintEntireIndex(pages []*spider.Page) {
	// from http://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file

	fo, err := os.Create("spider_result.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	for i, page := range pages {

		fo.WriteString(page.Title + "\n" + page.URL)
		fo.WriteString("\n")
		fo.WriteString(page.Modified.String() + ", ")
		fo.WriteString(strconv.FormatInt(page.Size, 10))
		fo.WriteString("\n")

		for word, wordFreq := range page.words {
			fo.WriteString(word + ", " + strconv.Itoa(len(wordFreq)) + "; ")
		}
		fo.WriteString("\n")

		//pageLinks := strings.Fields(page.childLinks)
		for _, link := range page.Links {
			fo.WriteString(link + "\n")
		}

		fo.WriteString("-------------------------------------------------------------------------------------------\n")
	}

}

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

	//table := indexHandler
	//table.InitialIndexAndMaps()
	//SetTablesFromDB(table)

	// pages := getPagesFromDb()
	PrintEntrieIndex(spider.Get30Pages())
}
