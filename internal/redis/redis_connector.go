package redis

import (
	"context"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type RedisCfg struct {
	Addr     string `json:"address" mapstructure:"address"`
	Password string `json:"password" mapstructure:"password"`
	DB       int    `json:"db" mapstructure:"db"`
	PoolSize int    `json:"pool_size" mapstructure:"pool_size"`
}

type RedisFailOverCfg struct {
	DB              int      `json:"db" mapstructure:"db"`
	Master          string   `json:"master" mapstructure:"master"`
	PoolSize        int      `json:"pool_size" mapstructure:"pool_size"`
	SentinelServers []string `json:"sentinels" mapstructure:"sentinels"`
	Password        string   `json:"password" mapstructure:"password"`
}

type RedisService interface {
	Set(context.Context, string, interface{}, time.Duration) error
	SetIfNotExist(context.Context, string, interface{}, time.Duration) (bool, error)
	Get(context.Context, string) *redis.StringCmd
	Delete(context.Context, string) error
	Client() *redis.Client
}

type redisService struct {
	failOverCfg *RedisFailOverCfg
	Cfg         *RedisCfg
	redisClient *redis.Client
	redisHosts  string
	redisPort   string
	redisDB     string
}

func NewFailOverClient(cfg *RedisFailOverCfg) (RedisService, error) {
	redisService := &redisService{
		failOverCfg: cfg,
	}
	redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.Master,
		SentinelAddrs: cfg.SentinelServers,
		Password:      cfg.Password,
		PoolSize:      cfg.PoolSize,
		DB:            cfg.DB,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	redisService.redisClient = redisClient
	redisService.getRedisHostsAndPort()
	return redisService, nil
}

func New(cfg *RedisCfg) (RedisService, error) {
	redisService := &redisService{
		Cfg: cfg,
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		PoolSize: cfg.PoolSize,
		DB:       cfg.DB,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	redisService.redisClient = redisClient

	hostAndPort := strings.Split(cfg.Addr, ":")
	if len(hostAndPort) > 1 {
		host, port := hostAndPort[0], hostAndPort[1]
		redisService.redisHosts = host
		redisService.redisPort = port
	}
	return redisService, nil
}

func (this *redisService) getRedisHostsAndPort() (hosts, port string) {
	for _, val := range this.failOverCfg.SentinelServers {
		host := strings.Split(val, ":")[0]
		hosts = hosts + host + ","
	}
	if len(hosts) == 0 {
		return
	}

	hosts = hosts[:len(hosts)-1]
	port = strings.Split(this.failOverCfg.SentinelServers[0], ":")[1]
	this.redisHosts = hosts
	this.redisPort = this.redisPort
	return
}

func (this *redisService) Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	err := this.redisClient.Set(key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

func (this *redisService) Get(ctx context.Context, key string) *redis.StringCmd {
	return this.redisClient.Get(key)
}

func (this *redisService) Delete(ctx context.Context, key string) error {
	_, err := this.redisClient.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (this *redisService) SetIfNotExist(ctx context.Context, key string, value interface{}, expiry time.Duration) (bool, error) {
	isSet, err := this.redisClient.SetNX(key, value, expiry).Result()
	if err != nil {
		return isSet, err
	}

	return isSet, err
}

func (this *redisService) Client() *redis.Client {
	return this.redisClient
}
