package cache

import (
	"banners/internal/database/models"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, expires time.Duration) BannerCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: expires,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value *models.Banner) {
	client := cache.getClient()

	banner, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	client.Set(key, banner, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) *models.Banner {
	client := cache.getClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	banner := models.Banner{}

	err = json.Unmarshal([]byte(val), &banner)
	if err != nil {
		panic(err)
	}

	return &banner
}
