package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"time"
)


// mockery -name=CacheServiceInterface --case underscore --inpkg true
type CacheServiceInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiry int) error
	GetGet(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	SetIfNotExist(ctx context.Context, key string, value interface{}, expiry int) (bool, error)
}

type cacheService struct {
	config *config.GeneralConfig
	redisClient *redis.Client
}

func (c *cacheService) Set(ctx context.Context, key string, value interface{}, expiry int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.redisClient.Set(key, string(data), time.Duration(expiry)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheService) Get(ctx context.Context, key string) (string, error) {
	val, err := c.redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *cacheService) Delete(ctx context.Context, key string) error {
	_, err := c.redisClient.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheService) SetIfNotExist(ctx context.Context, key string, value interface{}, expiry int) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	isSet, err := c.redisClient.SetNX(key, string(data), time.Duration(expiry)*time.Second).Result()
	if err != nil {
		return isSet, err
	}

	return isSet, err
}
