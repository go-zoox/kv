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
	test.RunTestCases(t, createClient())
}
