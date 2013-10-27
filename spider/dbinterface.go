package spider

// import (
// 	"fmt"
// 	"github.com/cznic/exp/dbm"
// 	"io"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// //const dbname = "4321.db"
// var o = &dbm.Options{}

// func SetTablesFromDB(table *IndexHandler) {
// 	var db *dbm.DB = dbConnect(dbname)
// 	/*invertedtable, _:= db.Array("invertedtable")
// 	enum,err := invertedtable.Enumerator(true)
// 	if err!=io.EOF{

// 	} else{
// 		key, value,err :=enum.Next()
// 		for ;err!=io.EOF; key, value,err =enum.Next(){
// 			wid:=key[0].(int64)

// 			var token []string = Fields(value[0].(string))
// 			for num,_:range token{
// 				pid:=token[num]
// 				table.invertedIndex[wid] = pid
// 			}
// 		}
// 	}*/
// 	fowardtable, _ := db.Array("fowardtable")
// 	enum, err := fowardtable.Enumerator(true)
// 	if err != io.EOF {

// 	} else {
// 		key, value, err := enum.Next()
// 		for ; err != io.EOF; key, value, err = enum.Next() {
// 			pid := key[0].(int64)

// 			token := strings.Split(value[0].(string), ";")
// 			for _, j := range token {

// 				//i:=0
// 				table.forwardIndex[int(pid)] = make(map[int][]int)
// 				var token2 []string = strings.Fields(j)
// 				positions := make([]int, 0, DefaultPositionsLength)
// 				wid := 0
// 				for i, temp := range token2 {
// 					//temp:=token[num]
// 					if i == 0 {
// 						wid, _ = strconv.Atoi(temp)
// 					} else {
// 						positionInt, _ := strconv.Atoi(temp)
// 						positions = append(positions, positionInt)
// 					}
// 				}
// 				table.forwardIndex[int(pid)][wid] = positions
// 			}
// 		}
// 	}

// 	for pid, wordIDs := range table.forwardIndex {

// 		for wid, _ := range wordIDs {
// 			_, exists := table.invertedIndex[wid]
// 			if !exists {
// 				table.invertedIndex[wid] = make(map[int]bool)
// 			}

// 			table.invertedIndex[wid][pid] = true
// 		}
// 	}

// 	GetWordMap(table)

// }

// func checkDbExist(filename string) bool {
// 	if _, err := os.Stat(filename); os.IsNotExist(err) {
// 		//fmt.Printf("no such file or directory: %s", filename)
// 		return false
// 	} else {
// 		return true
// 	}
// }

// func dbConnect(name string) *dbm.DB {
// 	if checkDbExist(name) != true {
// 		db, err := dbm.Create(name, o)
// 		if err != nil {
// 			fmt.Printf("Error: dbm can't create\n")
// 		} else {
// 			return db
// 		}
// 	} else {
// 		db, err := dbm.Open(name, o)
// 		if err != nil {
// 			fmt.Printf("Error: dbm can't open\n")
// 		} else {
// 			return db
// 		}
// 	}
// 	return nil
// }
