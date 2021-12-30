package db

import (
	"fmt"
	"os"

	"github.com/hskim881028/goblockchain/utility"
	bolt "go.etcd.io/bbolt"
)

const (
	// dbName       = "blockchain.db"
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDbName(), 0060, nil)
		db = dbPointer
		utility.HandleError(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utility.HandleError(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utility.HandleError(err)
	}
	return db
}

func PutBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utility.HandleError(err)
}

func GetBlock(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func DeleteAllBlocks() {
	DB().Update(func(t *bolt.Tx) error {
		utility.HandleError(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utility.HandleError(err)
		return nil
	})
}

func PutChain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utility.HandleError(err)
}

func GetChain() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Close() {
	DB().Close()
}

//for test
func getDbName() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}
