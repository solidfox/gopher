package spider

import (
	"fmt"
	"testing"
)

func TestIndexHandling(t *testing.T) {

	test := indexHandler
	test.StorePages(Get30Pages())
	test.PrintEntireIndex()

	//fmt.Print(Get30Pages())
	t.Fail()
}
