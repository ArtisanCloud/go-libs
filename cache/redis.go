package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	fmt2 "github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/redis/go-redis/v9"
)

type GRedis struct {
	Pool              redis.UniversalClient
	defaultExpiration time.Duration
	lockRetries       int
}

const SYSTEM_CACHE_TIMEOUT = 60 * 60
const SYSTEM_CACHE_TIMEOUT_MINUTE = 60
const SYSTEM_CACHE_TIMEOUT_HOUR = 60 * 60
const SYSTEM_CACHE_TIMEOUT_DAY = 60 * 60 * 24
const SYSTEM_CACHE_TIMEOUT_MONTH = 60 * 60 * 24 * 30
const SYSTEM_CACHE_TIMEOUT_SEASON = 60 * 60 * 24 * 30 * 3
const SYSTEM_CACHE_TIMEOUT_YEAR = 60 * 60 * 24 * 30 * 3 * 12

const (
	defaultMaxIdle            = 5
	defaultMaxActive          = 0
	defaultTimeoutIdle        = 240
	defaultTimeoutConnect     = 10000
	defaultTimeoutRead        = 5000
	defaultTimeoutWrite       = 5000
	defaultAddr               = "localhost:6379"
	defaultProtocol           = "tcp"
	defaultRetryThreshold     = 5
	defaultIdleCheckFrequency = time.Minute
)

const SCRIPT_SETEX = `return redis.call('exists',KEYS[1])<1 and redis.call('setex',KEYS[1],ARGV[2],ARGV[1])`

var CTXRedis = context.Background()

const lockRetries = 5

func NewGRedis(opts *redis.UniversalOptions) (gr *GRedis) {

	c := redis.NewUniversalClient(opts)
	gr = &GRedis{
		Pool:        c,
		lockRetries: lockRetries,
	}

	return gr

}

func (gr *GRedis) AddNX(key string, value interface{}, ttl time.Duration) bool {
	cmd := gr.Pool.SetNX(CTXRedis, key, value, ttl)
	r, err := cmd.Result()
	//fmt.Printf("r:%b \r\n", r)
	if err != nil {
		fmt.Errorf("SetNX error: %+v \r\n", err)
	}
	return r
}
func (gr *GRedis) Add(key string, value interface{}, ttl time.Duration) (err error) {
	// If the store has an "add" method we will call the method on the store so it
	// has a chance to override this logic. Some drivers better support the way
	// this operation should work with a total "atomic" implementation of it.
	var obj interface{}
	obj, err = gr.Get(key, obj)
	if err == ErrCacheMiss {
		return gr.SetEx(key, value, ttl)
	} else {
		return errors.New("this value has been actually added to the cache")
	}

}

func (gr *GRedis) Set(key string, value interface{}, expires time.Duration) error {
	mValue, err := json.Marshal(value)
	//mExpire, err := json.Marshal(expires)
	if err != nil {
		return err
	}

	result := gr.Pool.Set(CTXRedis, key, mValue, expires)
	return result.Err()
}

func (gr *GRedis) SetEx(key string, value interface{}, expires time.Duration) error {
	mValue, err := json.Marshal(value)
	//mExpire, err := json.Marshal(expires)
	if err != nil {
		return err
	}

	//luaScript := redis.NewScript(SCRIPT_SETEX)
	//cmd := luaScript.Run(gr.Pool.Context(),gr.Pool, []string{key}, mExpire, mValue)
	//fmt.Printf("result:%s \r\n",cmd.String())
	//fmt.Printf("err:%s \r\n", cmd.Err())

	// connPool := gr.Pool.Conn()
	cmd := gr.Pool.SetEx(CTXRedis, key, mValue, expires)
	//fmt2.Dump(connPool.Pipeline())
	//fmt.Printf("result:", cmd.String())

	return cmd.Err()
}

func (gr *GRedis) Get(key string, defaultValue interface{}) (ptrValue interface{}, err error) {
	b, err := gr.Pool.Get(CTXRedis, key).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &ptrValue)
	defaultValue = ptrValue
	return defaultValue, err
}

func (gr *GRedis) Has(key string) bool {

	value, err := gr.Get(key, nil)
	if value != nil && err == nil {
		return true
	}

	return false
}

func (gr *GRedis) GetMulti(keys ...string) (object.HashMap, error) {
	res, err := gr.Pool.MGet(CTXRedis, keys...).Result()
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrCacheMiss
	}

	m := make(object.HashMap)
	for ix, key := range keys {
		m[key] = res[ix].(string)
	}
	return m, nil
}

func (gr *GRedis) Delete(key string) error {
	return gr.Pool.Del(CTXRedis, key).Err()
}

func (gr *GRedis) Keys() ([]string, error) {
	return gr.Pool.Keys(CTXRedis, "*").Result()
}

func (gr *GRedis) Flush() error {
	return gr.Pool.FlushAll(CTXRedis).Err()
}

/**
 * Get an item from the cache, or execute the given Closure and store the result.
 *
 * @param  string  key
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @param  \Closure  callback
 * @return mixed
 */
func (gr *GRedis) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (obj interface{}, err error) {

	var value interface{}
	value, err = gr.Get(key, value)

	// If the item exists in the cache we will just return this immediately and if
	// not we will execute the given Closure and cache the result of that for a
	// given number of seconds so it's available for all subsequent requests.
	if err != nil && err != ErrCacheMiss {
		fmt2.Dump("error:", err.Error())
		return nil, err

	} else if value != nil {
		return value, err
	}

	value, err = callback()
	if err != nil {
		return nil, err
	}

	result := gr.Put(key, value, ttl)
	if !result {
		err = errors.New(fmt.Sprintf("remember cache put err, ttl:%d", ttl))
	}
	// ErrCacheMiss and query value from source
	return value, err
}

/**
 * Store an item in the cache.
 *
 * @param  string  key
 * @param  mixed  value
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @return bool
 */
func (gr *GRedis) Put(key interface{}, value interface{}, ttl time.Duration) bool {
	// key如果是数组
	//if arrayKey, ok := key.([]interface{}) !ok {
	//	return gr.PutMany(arrayKey, value)
	//}

	//if ttl == nil {
	//	return gr.forever(key, value)
	//}

	//seconds := gr.GetSeconds(ttl)
	//
	//if seconds <= 0 {
	//	return gr.Delete(key)
	//}

	//result = gr.Pool.Put(gr.itemKey(key), value, seconds)

	err := gr.SetEx(key.(string), value, ttl)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true

}

/**
 * Store multiple items in the cache for a given number of seconds.
 *
 * @param  array  values
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @return bool
 */
func (gr *GRedis) PutMany(values object.Array, ttl time.Duration) bool {
	//if ttl == nil {
	//	return gr.PutManyForever(values)
	//}
	//
	//seconds := gr.GetSeconds(ttl)
	//
	//if seconds <= 0 {
	//	return gr.Pool.Del(array_keys(values))
	//}
	//
	//gr.Pool.

	return false
}

/**
 * Store multiple items in the cache indefinitely.
 *
 * @param  array  values
 * @return bool
 */
func (gr *GRedis) PutManyForever(values []interface{}) bool {
	result := true

	//for key, value := range values {
	//
	//	if !gr.Forever(key, value) {
	//		result = false
	//	}
	//}

	return result
}

/**
 * Calculate the number of seconds for the given TTL.
 *
 * @param  \DateTimeInterface|\DateInterval|int  ttl
 * @return int
 */
func (gr *GRedis) GetSeconds(ttl time.Duration) int {
	//duration := gr.ParseDateInterval(ttl)
	//
	//if reflect.Type(duration).Kind() == DateTimeInterface {
	//	duration = carbon.Now().diffInRealSeconds(duration, false)
	//}
	//
	//if duration > 0 {
	//	return duration
	//} else {
	//	return 0
	//}
	return 0
}

func (gr *GRedis) SetByTags(key string, val interface{}, tags []string, expiry time.Duration) error {
	pipe := gr.Pool.TxPipeline()
	for _, tag := range tags {
		pipe.SAdd(CTXRedis, tag, key)
		pipe.Expire(CTXRedis, tag, expiry)
	}

	pipe.Set(CTXRedis, key, val, expiry)

	_, errExec := pipe.Exec(CTXRedis)
	return errExec
}

func (gr *GRedis) Invalidate(tags []string) {
	keys := make([]string, 0)
	for _, tag := range tags {
		k, _ := gr.Pool.SMembers(CTXRedis, tag).Result()
		keys = append(keys, tag)
		keys = append(keys, k...)
	}
	gr.Pool.Del(CTXRedis, keys...)
}
