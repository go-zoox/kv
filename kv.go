package kv

import (
	"github.com/go-zoox/kv/memory"
	"github.com/go-zoox/kv/redis"
	"github.com/go-zoox/kv/typing"
)

// NewMemory returns a new Memory KV.
func NewMemory() typing.KV {
	return memory.New()
}

// NewRedis returns a new Redis KV.
func NewRedis(cfg *redis.RedisConfig) (typing.KV, error) {
	return redis.New(cfg)
}
