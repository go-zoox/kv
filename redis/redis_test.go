package redis

import (
	"testing"

	"github.com/go-zoox/dotenv"
	"github.com/go-zoox/kv/test"
)

func createClient() *Redis {
	var cfg struct {
		URI string `env:"REDIS_URI"`
	}
	if err := dotenv.Load(&cfg); err != nil {
		panic(err)
	}

	client, err := New(&RedisConfig{
		URI:    cfg.URI,
		Prefix: "go-zoox-test:",
	})
	if err != nil {
		panic(err)
	}

	return client
}

func TestKV(t *testing.T) {
	client := createClient()
	client.Clear()
	defer client.Clear()

	test.RunMainTestCase(t, client)
	test.RunKeysTestCase(t, client)
	test.RunForEachTestCase(t, client)
}
