package spider

import (
	"fmt"
	"gopher/spider"
	"testing"
)

func TestParsePage(t *testing.T) {
	pages := spider.Get30Pages()
	for _, page := range pages {
		fmt.Println(page.URL)
	}
	t.Fail()
}
