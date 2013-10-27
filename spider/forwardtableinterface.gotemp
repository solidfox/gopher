package spider

import(
	"github.com/cznic/exp/dbm"
	"fmt"
	"io"
)

func StorePagesToForwardTable(table *IndexHandler) {
	// FWTable must not use URL as key but rather an ID.
	// Forward table should store TF, positions for words.
	// Should use ID -> URL table for meta data.
	
	var db *dbm.DB = dbConnect(dbname)
	fowardtable, _:= db.Array("fowardtable")
	// pid = pageid
	for pid,j1 := range table.forwardIndex{
		var temp string = ""
		for wid,j2 := range j1{
			temp = temp + fmt.Sprintf("%v",wid)
			for pos,_ := range j2{
				// f = position of the word
				temp = temp + " " + fmt.Sprintf("%v",pos)
			}
			temp = temp + ";"
		}
		fowardtable.Set(temp,pid)
	}
	closeDb(db)
}

/*func GetPagesFromDisk() (resultint []int64, resultStr []string) {
	//var resultint int[]
	var db *dbm.DB = dbConnect(dbname)
	fowardtable, _:= db.Array("fowardtable")
	enum,err := fowardtable.Enumerator(true)
	if err !=io.EOF{
	
	} else{
		//risky approach
		key, value,err :=enum.Next()
		for ;err!=io.EOF; key, value,err =enum.Next(){
			fmt.Printf("%v	%v\n\n",key[0].(int64),value[0].(string))
			resultint = append(resultint, key[0].(int64))
			resultStr = append(resultStr, value[0].(string))
		}
	}
	
	closeDb(db)
	return resultint,resultStr
}*/

func GetPagesFromForwardTable() ([]int64, []string){
	var resultint []int64
	var resultStr []string
	var db *dbm.DB = dbConnect(dbname)
	fowardtable, _:= db.Array("fowardtable")
	enum,err := fowardtable.Enumerator(true)
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

func DebugFowardtable(){
	var db *dbm.DB = dbConnect(dbname)
	fowardtable, _:= db.Array("fowardtable")
	enum,err := fowardtable.Enumerator(true)
	//risky approach
	key, value,err :=enum.Next()
	for ;err!=io.EOF; key, value,err =enum.Next(){
		
		fmt.Printf("%v	%v\n",key[0].(int64),value[0].(string))
	}
	closeDb(db)
	
}