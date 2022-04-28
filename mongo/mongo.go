package mongo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

// Mongo is a Key-Value Store in Mongo
type Mongo struct {
	sync.RWMutex
	Core   *goredis.Client
	Ctx    context.Context
	Config *MongoConfig
}

// MongoConfig is the configuration for Mongo
type MongoConfig struct {
	// Host is the host of the Mongo server
	Host string
	// Port is the port of the Mongo server
	Port int
	// Username is the username for the Mongo server
	Username string
	// Password is the password for the Mongo server
	Password string
	// Database is the database to use
	Database string

	// URI is the URI of the Mongo server
	// such as mongo://username:password@host:port/db
	URI string

	// Prefix is the prefix to use for all keys
	Prefix string
}

// New returns a new MemoryKV.
func New(cfg *MongoConfig) (*Mongo, error) {
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
	return &Mongo{
		Core:   core,
		Ctx:    ctx,
		Config: cfg,
	}, nil
}

func (m *Mongo) getKey(key string) string {
	return m.Config.Prefix + key
}

// Set sets the value for the given key.
// If maxAge is greater than 0, then the value will be expired after maxAge miliseconds.
func (m *Mongo) Set(key string, value string, maxAge ...int64) error {
	m.Lock()
	defer m.Unlock()

	maxAgeX := 0
	if len(maxAge) > 0 {
		maxAgeX = int(maxAge[0])
	}

	keyX := m.getKey(key)
	if maxAgeX > 0 {
		return m.Core.Set(m.Ctx, keyX, value, time.Duration(maxAgeX)*time.Millisecond).Err()
	}

	return m.Core.Set(m.Ctx, keyX, value, 0).Err()
}

// Get returns the value for the given key.
func (m *Mongo) Get(key string) string {
	m.RLock()
	defer m.RUnlock()

	keyX := m.getKey(key)
	return m.Core.Get(m.Ctx, keyX).Val()
}

// Delete deletes the value for the given key.
func (m *Mongo) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	return m.Core.Del(m.Ctx, m.getKey(key)).Err()
}

// Has returns true if the given key exists in the kv.
func (m *Mongo) Has(key string) bool {
	m.RLock()
	defer m.RUnlock()

	length, err := m.Core.Exists(m.Ctx, m.getKey(key)).Result()
	if err != nil {
		panic(err)
	}

	return length != 0
}

// Keys returns the keys of the kv.
func (m *Mongo) Keys() []string {
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
func (m *Mongo) Size() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.Keys())
}

// Clear removes all elements from the kv.
func (m *Mongo) Clear() error {
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
func (m *Mongo) ForEach(f func(string, interface{})) {
	m.RLock()
	defer m.RUnlock()

	keys := m.Keys()
	for _, key := range keys {
		f(key, m.Get(key))
	}
}
