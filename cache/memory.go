package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"os"
	"path"
	"time"
)

const DEFAULT_EXPIRES_IN = 5
const DEFAULT_PURGE_EXPIRES_IN_PERIOD = 10

type MemCache struct {
	*cache.Cache
	cacheFile string
}

func NewMemCache(namespace string, defaultLifeTime time.Duration, directory string) CacheInterface {

	if ACCache != nil {
		return ACCache
	}

	if defaultLifeTime <= 0 {
		defaultLifeTime = time.Duration(DEFAULT_EXPIRES_IN) * time.Minute
	}
	defaultPurgePeriod := time.Duration(DEFAULT_PURGE_EXPIRES_IN_PERIOD) * time.Minute

	path, err := createCacheFile(directory)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	memCache := &MemCache{
		cache.NewFrom(
			defaultLifeTime,
			defaultPurgePeriod,
			map[string]cache.Item{},
		),
		path,
	}

	//err = memCache.Cache.LoadFile(path)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	ACCache = memCache

	return ACCache
}

func createCacheFile(directory string) (cachePath string, err error) {

	_, err = os.Stat(directory)
	if err != nil && os.IsExist(err) {
		return "", err
	} else if os.IsNotExist(err) {
		directory, err = os.UserHomeDir()
	}

	directory = path.Join(directory, ".ArtisanCloud")
	err = os.Mkdir(directory, os.ModePerm)
	if err == nil || os.IsExist(err) {
		cachePath = path.Join(directory, "cache")
		_, err = os.Create(cachePath)
		if err == nil || os.IsExist(err) {
			return cachePath, nil
		}
	}

	if err != nil {
		return "", err
	}

	return cachePath, nil

}

func (cache *MemCache) Get(key string, defaultValue interface{}) (returnValue interface{}, err error) {
	ptrValue, found := cache.Cache.Get(key)
	if !found {
		return nil, errors.New(fmt.Sprintf("Cannot find value with key: %s", key))
	}

	err = json.Unmarshal(ptrValue.([]byte), &returnValue)
	return returnValue, err
}

func (cache *MemCache) Set(key string, value interface{}, expires time.Duration) error {

	mValue, err := json.Marshal(value)
	//mExpire, err := json.Marshal(expires)
	if err != nil {
		return err
	}

	cache.Cache.Set(key, mValue, expires)
	err = cache.Cache.SaveFile(cache.cacheFile)
	return err
}

func (cache *MemCache) Has(key string) bool {

	_, found := cache.Cache.Get(key)

	return found
}

func (cache *MemCache) AddNX(key string, value interface{}, ttl time.Duration) bool {
	return false
}

func (cache *MemCache) Add(key string, value interface{}, ttl time.Duration) (err error) {
	return nil
}

func (cache *MemCache) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (obj interface{}, err error) {
	return nil, nil
}
