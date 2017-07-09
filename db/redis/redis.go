package redis

/*
 * Redis client for temp storage and cache.
 */
import (
	redis "gopkg.in/redis.v3"
	"golangRestfulAPISample/bootstrap"
)

var redisClient *redis.Client

/*
 * Set connection to redis server
 */
func Init() {
	var err error
	redisClient = redis.NewClient(&redis.Options{
		Addr:     bootstrap.App.DBConfig.String("redis.address"),
		Password: bootstrap.App.DBConfig.String("redis.password"),
		DB:       0, // use default DB
	})
	if _, err = redisClient.Ping().Result();err != nil {
		panic(err)
	}
}

// GetValue : retrieve value from redis
func GetValue(key string) (interface{}, error) {
	var (
		val interface{}
		err error
	)

	if val, err = redisClient.Get(key).Result(); err != nil {
		return val, err
	}

	return val, nil
}

// SetValue : set value to redis
func SetValue(key string, value interface{}) error {
	var err error
	if err = redisClient.Set(key, value, 0).Err(); err != nil {
		return err
	}

	return nil
}

// DelKey : delete key from redis
func DelKey(key string) error {
	var err error
	if _, err = redisClient.Del(key).Result(); err != nil {
		return err
	}
	return nil
}
