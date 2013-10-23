package spider

import (
	"os"
	"strconv"
)

type indexHandler struct {
	// page		// word	// freq
	forwardIndex map[int]map[int][]int

	// word		// page
	invertedIndex map[int]map[int]bool

	wordMap        map[int]string
	wordMapReverse map[string]int

	// this page's
	pageMap        map[int]Page
	pageMapReverse map[string]int
}

func (table *indexHandler) StorePages(pages []Page) {

	table.initialIndexAndMaps()

	for _, eachPage := range pages {

		//1. add to mapping table
		existPage, pageid := table.addPageToMappingTable(eachPage)

		// if the page had been entered, quit
		// actually we should check the last modified date
		if existPage == true {
			continue
		}

		// point to
		table.forwardIndex[pageid] = make(map[int][]int)

		for word, positions := range eachPage.words {
			_, wordid := table.addWordToMappingTable(word)

			// assume this is by ref, and the memory won't be deleted
			table.forwardIndex[pageid][wordid] = positions

			_, exists := table.invertedIndex[wordid]
			if !exists {
				table.invertedIndex[wordid] = make(map[int]bool)
			}

			table.invertedIndex[wordid][pageid] = true

		}

	}
}

func (ih *indexHandler) PrintEntireIndex() {
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

	//buf := make([]byte, 1024)

	for pageid, eachPage := range ih.pageMap {

		fo.WriteString(eachPage.Title + "\n" + eachPage.URL)
		fo.WriteString("\n")
		fo.WriteString(eachPage.Modified.String() + ", ")
		fo.WriteString(strconv.FormatInt(eachPage.Size, 10))
		fo.WriteString("\n")

		for eachWordId, eachWordFreq := range ih.forwardIndex[pageid] {
			fo.WriteString(ih.wordMap[eachWordId] + ", " + strconv.Itoa(len(eachWordFreq)) + "; ")
		}

		fo.WriteString("\n")

		for _, eachLink := range ih.pageMap[pageid].links {
			fo.WriteString(ih.pageMap[ih.pageMapReverse[eachLink.URL]].Title + " " + eachLink.URL + "\n")
		}

		fo.WriteString("-------------------------------------------------------------------------------------------\n")

	}
}

func (table *indexHandler) initialIndexAndMaps() {

	table.wordMap = make(map[int]string)
	table.pageMap = make(map[int]Page)

	// will change to db if necessary
	table.forwardIndex = make(map[int]map[int][]int)
	table.invertedIndex = make(map[int]map[int]bool)
}

func (table *indexHandler) addPageToMappingTable(page Page) (bool, int) {
	_, ok := table.pageMapReverse[page.URL]
	if !ok {
		table.pageMap[len(table.pageMap)] = page
		table.pageMapReverse[page.URL] = len(table.pageMap) - 1

	}
	return ok, table.pageMapReverse[page.URL]
}

func (table *indexHandler) addWordToMappingTable(word string) (bool, int) {
	_, ok := table.wordMapReverse[word]
	if !ok {
		table.wordMap[len(table.wordMap)] = word
		table.wordMapReverse[word] = len(table.wordMap) - 1

	}
	return ok, table.wordMapReverse[word]
}
