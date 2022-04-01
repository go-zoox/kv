package redis

import (
	"testing"

	"github.com/go-zoox/dotenv"
	"github.com/go-zoox/kv/test"
)

func createClient() *Redis {
	// var cfg struct {
	// 	URI string `env:"REDIS_URI"`
	// }
	// if err := dotenv.Load(&cfg); err != nil {
	// 	panic(err)
	// }

	redisURI := dotenv.Get("REDIS_URI", "redis://localhost:6379")

	client, err := New(&RedisConfig{
		URI:    redisURI,
		Prefix: "go-zoox-test:",
	})
	if err != nil {
		panic(err)
	}

	return client
}

func TestKV(t *testing.T) {
	test.RunTestCases(t, createClient())
}
