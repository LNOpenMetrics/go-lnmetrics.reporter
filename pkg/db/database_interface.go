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

func (this *database) Ready() bool {
	return this.instance != nil
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

func (this *database) PutValue(key string, value string) error {
	log.GetInstance().Debug(
		fmt.Sprintf("Storing value with key %s and value %s", key, value))
	return this.instance.Put([]byte(key), []byte(value), nil)
}

func (this *database) GetValue(key string) (string, error) {
	log.GetInstance().Debug(fmt.Sprintf("Search on Db value with key %s", key))
	value, err := this.instance.Get([]byte(key), nil)
	log.GetInstance().Debug("Return value from db")
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("%s", err))
		return "", err
	}
	return string(value), nil
}

func (this *database) DeleteValue(key string) error {
	return this.instance.Delete([]byte(key), nil)
}

// TODO Add method to iterate over a method.
