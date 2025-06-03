package services

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
	"golang.org/x/net/context"
)

type IRedisService interface {
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
}

type RedisService struct {
	client redis.Cmdable
	ctx    context.Context
}

// NewRedisService creates and initializes a new Redis service connection
//
// Returns:
//   - *RedisService: Pointer to initialized Redis service
//   - Panics if connection fails
func NewRedisService(client redis.Cmdable) *RedisService {
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	} else {
		logger.Info("Connected to Redis successfully")
	}

	return &RedisService{
		client: client,
		ctx:    context.Background(),
	}
}

// Set stores a key-value pair in Redis with no expiration time
// Parameters:
//   - key: the key to store the value under
//   - value: the value to store (can be any type that Redis supports)
//
// Returns:
//   - error: nil if successful, otherwise contains the error message
func (r *RedisService) Set(key string, value any, ttl time.Duration) error {
	err := r.client.Set(r.ctx, key, value, ttl).Err()
	if err != nil {
		return errors.New(errors.ErrCacheSet, err.Error())
	}
	return nil
}

// Get retrieves a value from Redis by its key
// Parameters:
//   - key: the key to look up in Redis
//
// Returns:
//   - string: the value stored at the key, or empty string if key doesn't exist
//   - error: nil if successful or key not found, otherwise contains error message
func (r *RedisService) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", errors.New(errors.ErrCacheGet, err.Error())
	}
	return val, err
}

// Delete removes a key-value pair from Redis
// Parameters:
//   - key: the key to delete from Redis
//
// Returns:
//   - error: nil if successful, otherwise contains the error message
func (r *RedisService) Delete(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return errors.New(errors.ErrCacheDelete, err.Error())
	}
	return nil
}

// Exists checks if a key exists in Redis
// Parameters:
//   - key: the key to check in Redis
//
// Returns:
//   - bool: true if key exists, false if it doesn't
//   - error: nil if successful, otherwise contains the error message
func (r *RedisService) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, errors.New(errors.ErrCacheExists, err.Error())
	}
	return count > 0, nil
}
