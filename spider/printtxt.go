package spider

import (
	//"runtime"
	"os"
	"strconv"
)

func PrintEntireIndex(pages []*Page) {
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

	for _, page := range pages {

		fo.WriteString(page.Title)
		fo.WriteString("\n")
		fo.WriteString(page.URL)
		fo.WriteString("\n")
		fo.WriteString(page.Modified.String() + ", ")
		fo.WriteString(strconv.FormatInt(page.Size, 10))
		fo.WriteString("\n")

		for _, word := range page.Words() {
			fo.WriteString(word.Word + ", " + strconv.Itoa(len(word.Positions())) + "; ")
		}
		fo.WriteString("\n")

		//pageLinks := strings.Fields(page.childLinks)
		for _, link := range page.Links() {
			fo.WriteString(link.URL + "\n")
		}

		fo.WriteString("-------------------------------------------------------------------------------------------\n")
	}

}
