package kv

import (
	"github.com/go-zoox/kv/fs"
	"github.com/go-zoox/kv/memory"
	"github.com/go-zoox/kv/redis"

	"github.com/go-zoox/kv/typing"
)

// KV is the interface for all KV implementations.
type KV = typing.KV

// Config is the interface for KV Config.
type Config = typing.Config

// New returns a new KV.
func New(cfg *typing.Config) (KV, error) {
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

		return NewRedis(cfg.Config.(*redis.Config))
	default:
		return nil, NewError(ErrUnknownEngine, cfg.Engine)
	}
}

// NewMemory returns a new Memory KV.
func NewMemory() KV {
	return memory.New()
}

// NewFileSystem returns a new FileSystem KV.
func NewFileSystem(cfg ...*fs.FileSystemOptions) (KV, error) {
	return fs.New(cfg...)
}

// NewRedis returns a new Redis KV.
func NewRedis(cfg *redis.Config) (KV, error) {
	return redis.New(cfg)
}
