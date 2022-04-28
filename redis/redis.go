package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

// Redis is a Key-Value Store in Redis
type Redis struct {
	sync.RWMutex
	Core   *goredis.Client
	Ctx    context.Context
	Config *RedisConfig
}

// RedisConfig is the configuration for Redis
type RedisConfig struct {
	// Host is the host of the Redis server
	Host string
	// Port is the port of the Redis server
	Port int
	// Password is the password for the Redis server
	Password string
	// DB is the database to use
	DB int

	// URI is the URI of the Redis server
	// such as redis://:password@host:port/db
	URI string

	// Prefix is the prefix to use for all keys
	Prefix string
}

// New returns a new MemoryKV.
func New(cfg *RedisConfig) (*Redis, error) {
	var core *goredis.Client
	if cfg.URI != "" {
		opt, err := goredis.ParseURL(cfg.URI)
		if err != nil {
			return nil, err
		}
		core = goredis.NewClient(opt)
	} else if cfg.Host != "" && cfg.Port != 0 {
		core = goredis.NewClient(&goredis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})
	} else {
		return nil, fmt.Errorf("redis URI or Host and Port are required")
	}

	if cfg.Prefix == "" {
		return nil, errors.New("prefix is required")
	}

	// @TODO
	ctx := context.Background()
	return &Redis{
		Core:   core,
		Ctx:    ctx,
		Config: cfg,
	}, nil
}

func (m *Redis) getKey(key string) string {
	return m.Config.Prefix + key
}

func (m *Redis) encodeValue(value any) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func (m *Redis) decodeValue(data []byte, value any) error {
	return json.Unmarshal(data, value)
}

// Set sets the value for the given key.
// If maxAge is greater than 0, then the value will be expired after maxAge miliseconds.
func (m *Redis) Set(key string, value any, maxAge ...int64) error {
	m.Lock()
	defer m.Unlock()

	maxAgeX := 0
	if len(maxAge) > 0 {
		maxAgeX = int(maxAge[0])
	}

	keyX := m.getKey(key)
	valueX, err := m.encodeValue(value)
	if err != nil {
		return err
	}

	if maxAgeX > 0 {
		return m.Core.Set(m.Ctx, keyX, valueX, time.Duration(maxAgeX)*time.Millisecond).Err()
	}

	return m.Core.Set(m.Ctx, keyX, valueX, 0).Err()
}

// Get returns the value for the given key.
func (m *Redis) Get(key string, value any) error {
	m.RLock()
	defer m.RUnlock()

	keyX := m.getKey(key)
	valueX := m.Core.Get(m.Ctx, keyX).Val()
	return m.decodeValue([]byte(valueX), value)
}

// Delete deletes the value for the given key.
func (m *Redis) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	return m.Core.Del(m.Ctx, m.getKey(key)).Err()
}

// Has returns true if the given key exists in the kv.
func (m *Redis) Has(key string) bool {
	m.RLock()
	defer m.RUnlock()

	length, err := m.Core.Exists(m.Ctx, m.getKey(key)).Result()
	if err != nil {
		panic(err)
	}

	return length != 0
}

// Keys returns the keys of the kv.
func (m *Redis) Keys() []string {
	m.RLock()
	defer m.RUnlock()

	res := m.Core.Keys(m.Ctx, m.Config.Prefix+"*")
	if res.Err() != nil {
		panic(res.Err())
	}

	keys := make([]string, len(res.Val()))
	for i, k := range res.Val() {
		keys[i] = k[len(m.Config.Prefix):]
	}
	return keys
}

// Size returns the number of elements in the kv.
func (m *Redis) Size() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.Keys())
}

// Clear removes all elements from the kv.
func (m *Redis) Clear() error {
	keys := m.Keys()
	if len(keys) == 0 {
		return nil
	}

	m.Lock()
	defer m.Unlock()

	for _, key := range keys {
		// m.Delete(key)
		if err := m.Core.Del(m.Ctx, m.getKey(key)).Err(); err != nil {
			return err
		}
	}
	return nil
}

// ForEach calls the given function for each key-value pair in the kv.
func (m *Redis) ForEach(f func(string, interface{})) {
	m.RLock()
	defer m.RUnlock()

	keys := m.Keys()
	for _, key := range keys {
		var value any
		if err := m.Get(key, &value); err != nil {
			f(key, nil)
		} else {
			f(key, value)
		}
	}
}
