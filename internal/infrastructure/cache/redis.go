package cache

import (
	"encoding/json"
	"fmt"
	"shorter-rest-api/internal/config"
	"shorter-rest-api/internal/domain/entity"
	"time"

	"github.com/gomodule/redigo/redis"
)

type IRedisCache interface {
	Set(key string, value entity.ShortURL, expiration int) error
	CountKeysByPattern(pattern string) (int, error)
	Get(key string) (*entity.ShortURL, error)
	Exists(key string) (bool, error)
}

// RedisClient represents a Redis client
type RedisClient struct {
	Conn *redis.Pool
}

func (r *RedisClient) CountKeysByPattern(pattern string) (int, error) {
	var count int
	var cursor uint64 = 0

	conn := r.Conn.Get()
	defer conn.Close()
	for {
		keys, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", pattern, "COUNT", 100))
		if err != nil {
			return 0, err
		}

		if len(keys) != 2 {
			return 0, fmt.Errorf("invalid SCAN response")
		}

		// Extract cursor and keys
		cursor, _ = keys[0].(uint64)
		keySlice, _ := redis.Strings(keys[1], nil)
		count += len(keySlice)

		// If cursor is 0, we've scanned all keys
		if cursor == 0 {
			break
		}
	}

	return count, nil
}

// Set sets a key-value pair in Redis with expiration
func (r *RedisClient) Set(key string, value entity.ShortURL, expiration int) error {
	conn := r.Conn.Get()

	defer conn.Close()
	// Marshal the value to JSON
	rawData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	// If no expiration is set, use 0
	// which means the key will not expire
	return conn.Send("SET", key, rawData, "EX", expiration)
}

// Get gets a value from Redis by key
func (r *RedisClient) Get(key string) (*entity.ShortURL, error) {

	conn := r.Conn.Get()
	defer conn.Close()

	rawData, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("short_urls:%s", key)))
	if err != nil {
		return nil, fmt.Errorf("failed to get value from Redis: %w", err)
	}
	if len(rawData) <= 0 {
		return nil, nil
	}
	var shortUrl entity.ShortURL
	if err := json.Unmarshal(rawData, &shortUrl); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return &shortUrl, nil
}

func (r *RedisClient) Exists(key string) (bool, error) {
	conn := r.Conn.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, fmt.Errorf("failed to check if key exists: %w", err)
	}
	defer conn.Close()

	return exists, nil
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.Config) (IRedisCache, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	redisPool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, redis.DialPassword(cfg.Redis.Password))
		},
	}
	return &RedisClient{Conn: redisPool}, nil
}
