package db

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"

	log "github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

type database struct {
	instance *leveldb.DB
}

var instance database

func GetInstance() *database {
	return &instance
}

func (this *database) InitDB(homedir string) error {
	path := homedir + "/db"
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.GetInstance().Error(err)
		return err
	}
	log.GetInstance().Info("Created database at " + path)
	this.instance = db
	return nil
}

func (this *database) PutValue(key string, value interface{}) error {
	return this.instance.Put([]byte(key), []byte(fmt.Sprintf("%v", value)), nil)
}

func (this *database) GetValue(key string) (interface{}, error) {
	return this.instance.Get([]byte(key), nil)
}

func (this *database) DeleteValue(key string) error {
	return this.instance.Delete([]byte(key), nil)
}

// TODO Add method to iterate over a method.
