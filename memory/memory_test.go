package memory

import (
	"testing"

	"github.com/go-zoox/kv/test"
)

func createClient() *Memory {
	client := New()
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
