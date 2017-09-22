// ethdb数据处理
package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	_ "github.com/syndtr/goleveldb/leveldb"
)

func DbPut(key, value string, ldb *ethdb.LDBDatabase) {
	if err := ldb.Put([]byte(key), []byte(value)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Put OK")
}

func DbGet(key string, ldb *ethdb.LDBDatabase) {
	val, err := ldb.Get([]byte(key)); 
	if err != nil {
		fmt.Println("DB Get : ", err)
		return
	}
	fmt.Println("GetData: ", val)
}

func main() {
	ldb, err := ethdb.NewLDBDatabase("test.db", 10,  10)
	if err != nil {
		fmt.Println(err)
		return 
	}

	defer ldb.Close()

	DbPut("Name", "Bill", ldb)
	DbGet("Name", ldb)
}