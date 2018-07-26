package storage

import (
	"github.com/go-redis/redis"
)

type RedisStorage struct {
	redis *redis.Client
}

type RedisOption struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	PoolSize int    `json:"pool_size""`
}

func CreateRedisStorage(option *RedisOption) (*RedisStorage, error) {
	o := &redis.Options{
		Addr:     option.Addr,
		Password: option.Password,
	}
	if option.PoolSize > 0 {
		o.PoolSize = option.PoolSize
	}
	client := redis.NewClient(o)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &RedisStorage{
		redis: client,
	}, nil
}

func (rs *RedisStorage) DB() *redis.Client {
	return rs.redis
}

func CreateMultiRedisStorage(options map[string]*RedisOption) (map[string]*RedisStorage, error) {
	result := make(map[string]*RedisStorage)
	for name, option := range options {
		o := &redis.Options{
			Addr:     option.Addr,
			Password: option.Password,
		}
		if option.PoolSize > 0 {
			o.PoolSize = option.PoolSize
		}
		client := redis.NewClient(o)
		_, err := client.Ping().Result()
		if err != nil {
			return nil, err
		}
		result[name] = &RedisStorage{
			redis: client,
		}
	}
	return result, nil
}
