package spider

// import (
// 	"fmt"
// 	"github.com/cznic/exp/dbm"
// 	"io"
// )

// func StorePagesToPageInfo(table *IndexHandler) {

// 	for wid, w := range table.wordMap {

// 		SetWordId(w, wid)

// 	}

// 	pageInfoDB := NewRelationalDB("sqlite.db")

// 	defer pageInfoDB.Close()

// 	for _, page := range table.pageMap {
// 		var childLinks string = ""

// 		links := page.Links()
// 		for _, link := range links {
// 			childLinks = childLinks + link.URL + " "
// 		}

// 		pageInfo := PageInfo{
// 			int64(table.pageMapReverse[page.URL]),
// 			page.Size,
// 			page.URL,
// 			page.Modified,
// 			page.Title,
// 			childLinks,
// 		}

// 		pageInfoDB.InsertPageInfo(&pageInfo)
// 	}
// }

// func DebugPageInfo() {
// 	var db *dbm.DB = dbConnect(dbname)
// 	invertedtable, _ := db.Array("invertedtable")
// 	enum, err := invertedtable.Enumerator(true)
// 	//risky approach
// 	key, value, err := enum.Next()
// 	for ; err != io.EOF; key, value, err = enum.Next() {
// 		fmt.Printf("%v	%v\n", key[0].(int64), key[1].(string), value[0].(string))
// 	}
// 	db.Close()
// }
