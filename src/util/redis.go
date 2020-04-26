package util

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

func formatRedisAddress(addr string, port int) string {
	return fmt.Sprintf("%s:%d", addr, port)
}

func GetSessionRedisClient() *redis.Client {
	config := GetConfig().Redis
	client := redis.NewClient(&redis.Options{
		Addr:     formatRedisAddress(config.Address, config.Port),
		Password: config.Password,
		DB:       config.Session,
	})

	return client
}

func GetRegisterRedisClient() *redis.Client {
	config := GetConfig().Redis
	client := redis.NewClient(&redis.Options{
		Addr:     formatRedisAddress(config.Address, config.Port),
		Password: config.Password,
		DB:       config.Register,
	})

	return client
}

func GetUploadRedisClient() *redis.Client {
	config := GetConfig().Redis
	client := redis.NewClient(&redis.Options{
		Addr:     formatRedisAddress(config.Address, config.Port),
		Password: config.Password,
		DB:       config.Upload,
	})
	return client
}
