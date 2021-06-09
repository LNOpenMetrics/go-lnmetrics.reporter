package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Database struct {
	instance *DB
}

func (this *Database) InitDB(homedir String) err {
	db, err := leveldb.OpenFile("path/to/db", nil)
	if err != nil {
		return err
	}
	this.instance = db
	return nil
}

func (this *Database) PutValue(key String, value *interface{}) err {
	return this.instance.Put([]byte(key), []byte(valye), nil)
}

func (this *Database) GetValue(key String) (interface{}, err) {
	return this.instance.Get([]byte(key), []byte(), nil)
}

func (this *Database) DeleteValue(key String) err {
	return this.instance.Delete([]byte(key), nil)
}

// TODO Add method to iterate over a method.
