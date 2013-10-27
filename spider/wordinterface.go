package spider

// import(
// 	"github.com/cznic/exp/dbm"
// 	"fmt"
// 	"io"
// )

// //const dbname = "4321.db"

// func GetWordId(URL string){
// 	var db *dbm.DB = dbConnect(dbname)
// 	wordtoPageId, _:= db.Array("wordtoPageId")
// 	wid,err := wordtoPageId.Get(URL)

// }

// func GetWordMap(table *IndexHandler){
// 	var db *dbm.DB = dbConnect(dbname)
// 	wordtoPageId, _:= db.Array("wordtoPageId")
// 	enum,err := wordtoPageId.Enumerator(true)
// 	for ;err!=io.EOF; key, value,err =enum.Next(){
// 		table.wordMap[value[0].(int64)] = key[0].(string)

// 	}
// }

// func generateWId(a dbm.Array) int64{
// 	enum,err := a.Enumerator(false)
// 	//risky approach
// 	if err !=nil{
// 		return 0
// 	}
// 	//key, value, err := enum.Next()
// 	_, value, _ := enum.Next()
// 	i:=value[0].(int64)+1
// 	/*key, value, err := enum.Next()*/
// 	//must corr approach
// 	i:=1
// 	for ;err!=io.EOF; key2, value,err :=enum.Next(){

// 		i++
// 	}
// 	return i
// }

// func SetWordId(URL string, wid int) {
// 	var db *dbm.DB = dbConnect(dbname)
// 	wordtoPageId, _:= db.Array("wordtoPageId")
// 	//wid,_ := wordtoPageId.Get(URL)

// 	wordtoPageId.Set(wid,URL)
// 	defer closeDb(db)

// /*	if wid == nil{
// 		i := generateWId(wordtoPageId)
// 		//var i int64 = 0
// 		wordtoPageId.Set(i,URL)
// 		return i
// 	} else{
// 		return wid.(int64)
// 	}
// */
// //	return -1
// }

// /*func closeDb(db *dbm.DB){
// 	err2 := db.Close()
// 	if err2 != nil{
// 		fmt.Printf("can't close")
// 	}
// }*/

// func DebugWid(){
// 	var db *dbm.DB = dbConnect(dbname)
// 	a, _:= db.Array("wordtoPageId")
// 	enum,err := a.Enumerator(true)
// 	//risky approach
// key2, value,err :=enum.Next()
// 	for ;err!=io.EOF; key2, value,err =enum.Next(){

// 		fmt.Printf("%v	%v\n",key2[0].(string),value[0].(int64))
// 	}
// 	db.Close()
// }
