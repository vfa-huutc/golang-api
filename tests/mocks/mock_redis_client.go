package mocks

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

// --- Mock Redis Client ---
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, ttl)
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	} else {
		cmd.SetVal(args.String(1))
	}
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewIntCmd(ctx)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewIntCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	} else {
		cmd.SetVal(int64(args.Int(1)))
	}
	return cmd
}

func (m *MockRedisClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	args := m.Called(ctx, key, start, stop)
	cmd := redis.NewStringSliceCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	} else {
		cmd.SetVal(args.Get(1).([]string))
	}
	return cmd
}
