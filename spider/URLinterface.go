package spider

import(
	"github.com/cznic/exp/dbm"
	"fmt"
	"io"
)

const dbname = "4321.db"

func GetPageId(URL string) int64{
	var db *dbm.DB = dbConnect(dbname)
	URLtoPageId, _:= db.Array("URLtoPageId")
	pid,err := URLtoPageId.Get(URL)
	if err != nil{
		closeDb(db)
		return 0
	} else{
		closeDb(db)
		return pid.(int64)
	}
	
	return -1
}

func generatePId(a dbm.Array) int64{
	enum,err := a.Enumerator(false)
	//risky approach
	if err !=nil{
		return 0
	}
	//key, value, err := enum.Next()
	_, value, _ := enum.Next()
	i:=value[0].(int64)+1
	/*key, value, err := enum.Next()*/
	/*//must corr approach
	i:=1
	for ;err!=io.EOF; key2, value,err :=enum.Next(){
		
		i++
	}*/
	return i
}

func SetPageId(URL string) int64{
	var db *dbm.DB = dbConnect(dbname)
	URLtoPageId, _:= db.Array("URLtoPageId")
	pid,_ := URLtoPageId.Get(URL)
	if pid == nil{
		i := generatePId(URLtoPageId)
		//var i int64 = 0
		URLtoPageId.Set(i,URL)
		closeDb(db)
		return i
	} else{
		closeDb(db)
		return pid.(int64)
	}
	
	return -1
}

func closeDb(db *dbm.DB){
	err2 := db.Close()
	if err2 != nil{
		fmt.Printf("can't close")
	}
}

func DebugPid(){
	var db *dbm.DB = dbConnect(dbname)
	a, _:= db.Array("URLtoPageId")
	enum,err := a.Enumerator(true)
	//risky approach
key2, value,err :=enum.Next()
	for ;err!=io.EOF; key2, value,err =enum.Next(){
		
		fmt.Printf("%v	%v\n",key2[0].(string),value[0].(int64))
	}
	db.Close()
}