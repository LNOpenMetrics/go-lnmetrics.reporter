package cache

import (
	"encoding/json"
	"fmt"
	db "github.com/LNOpenMetrics/lnmetrics.utils/db/leveldb"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
	"strings"
)

// cacheManager is the internal rappresentation
// of the cache manager, that contains all the utils
// function to manage the cache
type cacheManager struct {
	prefix   string
	cacheIdx string
	cache    map[string]*string
}

// GetInstance return the instance of the cache manager that it is
// a singleton.
func GetInstance() *cacheManager {
	return &cacheManager{
		prefix:   "cache",
		cacheIdx: "cache/idx",
		cache:    nil,
	}
}

func (instance *cacheManager) buildID(key string) string {
	return strings.Join([]string{instance.prefix, key}, "/")
}

func (instance *cacheManager) initCache() error {
	instance.cache = make(map[string]*string)
	arrJson, err := json.Marshal(instance.cache)
	if err != nil {
		return err
	}
	if err := db.GetInstance().PutValueInBytes(instance.cacheIdx, arrJson); err != nil {
		log.GetInstance().Errorf("%s", err)
		return err
	}
	return nil
}

func (instance *cacheManager) getCacheIndex() (map[string]*string, error) {
	if instance.cache == nil {
		if err := instance.initCache(); err != nil {
			return nil, err
		}
		return instance.getCacheIndex()
	}
	value, err := db.GetInstance().GetValueInBytes(instance.cacheIdx)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(value, &instance.cache); err != nil {
		return nil, err
	}
	return instance.cache, nil
}

func (instance *cacheManager) addToCache(key string) error {
	idx := instance.cache
	if idx == nil {
		tmpIdx, err := instance.getCacheIndex()
		if err != nil {
			return err
		}
		idx = tmpIdx
	}
	if _, ok := idx[key]; !ok {
		idx[key] = &key
	}
	return nil
}

// IsInCache check if the key is inside the cache and return the result
// this not include side effect, so if any side effect happens, the function
// always return false.
func (instance *cacheManager) IsInCache(key string) bool {
	if instance.cache != nil {
		_, ok := instance.cache[key]
		if ok {
			return ok
		}
		// otherwise, continue and check with the database
	}
	key = instance.buildID(key)
	_, err := db.GetInstance().GetValue(key)
	if err != nil {
		log.GetInstance().Errorf("Error inside the cache: %s", err)
		return false
	}
	return true
}

// GetFromCache  retrieval the information that are in the cache in bytes
func (instance *cacheManager) GetFromCache(key string) ([]byte, error) {
	if instance.IsInCache(key) {
		key = instance.buildID(key)
		return db.GetInstance().GetValueInBytes(key)
	}
	return nil, fmt.Errorf("no value with key %s in the cache", key)
}

// PutToCache put the value in the json form inside the cache with the key specified
func (instance *cacheManager) PutToCache(key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	key = instance.buildID(key)
	if err := db.GetInstance().PutValueInBytes(key, jsonValue); err != nil {
		return err
	}
	return instance.addToCache(key)
}

// PurgeFromCache delete the value from cache with the specified key
// otherwise return an error.
func (instance *cacheManager) PurgeFromCache(key string) error {
	key = instance.buildID(key)
	if err := db.GetInstance().DeleteValue(key); err != nil {
		return err
	}
	return nil
}

// CleanCache clean the index from the database
func (instance *cacheManager) CleanCache() error {
	for key := range instance.cache {
		if err := db.GetInstance().DeleteValue(key); err != nil {
			return err
		}
	}
	return nil
}
