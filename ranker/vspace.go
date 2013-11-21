package ranker

import (
	//"fmt"
	"gopher/spider"
	"math"
)

func main() {
	db := spider.NewDBM("DBM.db")
	pages2 := db.GetPages2()

	db.Close()

}

func ifidf(documentID int, wordID int) (score Double) {
	db := spider.NewDBM("DBM.db")

	db.Close()
}
