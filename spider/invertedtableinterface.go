package spider

import (
	"github.com/cznic/exp/dbm"
	"fmt"
	"io"
)

func StorePagesToInvTable(table *IndexHandler) {
	// IVTable must not use string as key but rather an ID.
	var db *dbm.DB = dbConnect(dbname)
	invertedtable, _:= db.Array("invertedtable")
	// pid = pageid
	for wid,m1 := range table.invertedIndex{
		var temp string = ""
		for pid,_ := range m1{
			temp = temp + fmt.Sprintf("%v",pid)
		}
		invertedtable.Set(temp,wid)
	}
	closeDb(db)
}

func GetPagesFromInvertedTable() ([]int64, []string){
	var resultint []int64
	var resultStr []string
	var db *dbm.DB = dbConnect(dbname)
	invertedtable, _:= db.Array("invertedtable")
	enum,err := invertedtable.Enumerator(true)
	//risky approach
	key, value,err :=enum.Next()
	for ;err!=io.EOF; key, value,err =enum.Next(){
		
		fmt.Printf("%v	%v\n",key[0].(int64),value[0].(string))
		resultint = append(resultint, key[0].(int64))
		resultStr = append(resultStr, value[0].(string))
	}
	
	closeDb(db)
	return resultint,resultStr
}


func DebugInvertedtable(){
	var db *dbm.DB = dbConnect(dbname)
	invertedtable, _:= db.Array("invertedtable")
	enum,err := invertedtable.Enumerator(true)
	//risky approach
	key, value,err :=enum.Next()
	for ;err!=io.EOF; key, value,err =enum.Next(){
		fmt.Printf("%v	%v\n",key[0].(int64),value[0].(string))
	}
	db.Close()
	
}