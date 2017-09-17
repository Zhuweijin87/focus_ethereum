package main

/* 基于go的leveldb， 需要预先下载好leveldb包 */

import "fmt"
import "github.com/syndtr/goleveldb/leveldb"

func leveldb_put(db *leveldb.DB) {
	for i:=0; i<10; i++ {
		db.Put([]byte(fmt.Sprintf("key-%d", i)), []byte(fmt.Sprintf("A-Value-%02d", i)), nil)
	}
}

func leveldb_iter(db *leveldb.DB) {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println("Key: ", string(key), " Value:", string(value))
	}

	iter.Release()
	err := iter.Error()
	if err != nil {
		fmt.Println(err)
	}
}

func leveldb_get(key string, db *leveldb.DB) {
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func leveldb_getlike(key_like string, db *leveldb.DB) {

}

func main() {
	db, err := leveldb.OpenFile("./test.db", nil)
	if err != nil {
		fmt.Println("fail to open LevelDB:", err)
		return 
	}

	fmt.Printf("db type: %T\n", db)

	leveldb_put(db)

	leveldb_iter(db)

}