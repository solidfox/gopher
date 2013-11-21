package main

import (
	"fmt"
	"github.com/cznic/exp/dbm"

	"os"
)

var o = &dbm.Options{}

type DBM struct {
	db *dbm.DB
}

const DBMname = "DBM2.db"

func NewDBM(name string) (d *DBM) {
	mydb := dbConnect(name)

	d = &DBM{
		mydb,
	}
	return d
}
func checkDbExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		//fmt.Printf("no such file or directory: %s", filename)
		return false
	} else {
		return true
	}
}

func dbConnect(name string) *dbm.DB {
	if checkDbExist(name) != true {
		db, err := dbm.Create(name, o)
		if err != nil {
			fmt.Printf("Error: dbm can't create\n")
			panic(err)
		} else {
			return db
		}
	} else {
		db, err := dbm.Open(name, o)
		if err != nil {
			fmt.Printf("Error: dbm can't open\n")
			panic(err)
		} else {
			return db
		}
	}
	return nil
}
func main() {
	mydb := NewDBM(DBMname).db
	test, _ := mydb.Array("test")
	test.Set(1, "Peter", "Chan")
	test.Set(2, "Amy", "Lee")
	test.Set(3, "Dik", "Li")
	//test2, _ := mydb.Array("test", "Peter", "Chan")
	//test2, _ := test.Array("Peter", "Chan")
	test2, _ := mydb.Array("test", "Amy")
	val, _ := test2.Get("Lee")
	array, _ := dbm.MemArray("test")

	fmt.Printf("value=%v\n", val.(int64))
	enum, _ := array.Enumerator(true)
	_, value, _ := enum.Next()
	if value[0] == nil {
		fmt.Printf("Not OK")
	} else {
		fmt.Printf("%v", value[0].(int64))
	}

	// enum, _ := test2.Enumerator(true)
	// key, value, _ := enum.Next()
	// key, value, _ = enum.Next()
	// str := key[0].(string)
	// integer := value[0].(int64)
	// fmt.Printf("key=%v\nvalue=%v\n", str, integer)
	mydb.Close()

}
