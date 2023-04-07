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

func TestSetValueNil(t *testing.T) {
	client := createClient()
	defer client.Clear()
	if err := client.Set("key", nil); err == nil {
		t.Errorf("Expected error, got nil")
	}
}
