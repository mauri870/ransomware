package repository

import (
	"log"

	"github.com/tidwall/buntdb"
)

type Buntdb struct {
	*buntdb.DB
}

// Open or create a buntdb database
func Open(name string) *Buntdb {
	bunt, err := buntdb.Open(name)
	if err != nil {
		log.Fatalln(err)
	}
	return &Buntdb{bunt}
}

// Find by key
func (db Buntdb) Find(key string) (string, error) {
	value := ""
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		value = val
		return nil
	})
	return value, err
}

// Create or update a value
func (db *Buntdb) CreateOrUpdate(key string, value string) error {
	return db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
}

// Check if a key is available
func (db Buntdb) IsAvailable(searchKey string) bool {
	available := true
	db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(key, value string) bool {
			if searchKey == key {
				available = false
			}
			return false
		})
		return nil
	})

	return available
}
