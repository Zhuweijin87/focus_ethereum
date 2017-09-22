package main

/* 基于go的leveldb， 需要预先下载好leveldb包 */

import "fmt"
import "github.com/syndtr/goleveldb/leveldb"

/* insert */
func leveldb_put(db *leveldb.DB) {
	for i:=0; i<10; i++ {
		db.Put([]byte(fmt.Sprintf("key-%d", i)), []byte(fmt.Sprintf("A-Value-%02d", i)), nil)
	}
}

/* 迭代数据 */
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
	iter := db.NewIterator(util.BytesPrefix([]byte("key-")), nil)
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

/* 批量读写 */
func leveldb_batch_write(db *leveldb.DB) {
	batch := new(leveldb.Batch)
	batch.Put([]byte("key01"), []byte("value-01"))
	batch.Put([]byte("key02"), []byte("value-02"))
	batch.Put([]byte("key03"), []byte("value-03"))
	batch.Delete([]byte("key02"))

	err := db.Write(batch, nil)
	if err != nil {
		fmt.Println(err)
	}
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