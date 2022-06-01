package kv

import (
	"github.com/go-zoox/kv/fs"
	"github.com/go-zoox/kv/memory"
	"github.com/go-zoox/kv/redis"
	"github.com/go-zoox/kv/sqlite"
	"github.com/go-zoox/kv/typing"
)

// New returns a new KV.
func New(cfg *typing.Config) (typing.KV, error) {
	switch cfg.Engine {
	case "memory":
		return NewMemory(), nil

	case "filesystem":
		if cfg.Config == nil {
			return NewFileSystem()
		}

		return NewFileSystem(cfg.Config.(*fs.FileSystemOptions))

	case "redis":
		if cfg.Config == nil {
			return nil, NewError(ErrConfigNotSet, "redis")
		}

		return NewRedis(cfg.Config.(*redis.RedisConfig))
	case "sqlite":
		if cfg.Config == nil {
			return nil, NewError(ErrConfigNotSet, "sqlite")
		}

		return NewSQLite(cfg.Config.(*sqlite.SQLiteConfig))
	default:
		return nil, NewError(ErrUnknownEngine, cfg.Engine)
	}
}

// NewMemory returns a new Memory KV.
func NewMemory() typing.KV {
	return memory.New()
}

// NewFileSystem returns a new FileSystem KV.
func NewFileSystem(cfg ...*fs.FileSystemOptions) (typing.KV, error) {
	return fs.New(cfg...)
}

// NewRedis returns a new Redis KV.
func NewRedis(cfg *redis.RedisConfig) (typing.KV, error) {
	return redis.New(cfg)
}

// NewSQLite returns a new SQLite KV.
func NewSQLite(cfg *sqlite.SQLiteConfig) (typing.KV, error) {
	return sqlite.New(cfg)
}
