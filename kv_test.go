package kv

import (
	"testing"

	"github.com/go-zoox/kv/typing"
)

func TestKV(t *testing.T) {
	client, err := New(&typing.Config{
		// Engine: "sqlite",
		// Config: &sqlite.SQLiteConfig{
		// 	Path:   "/tmp/test.db",
		// 	Prefix: "go-zoox-test:",
		// },
		Engine: "memory",
	})
	if err != nil {
		t.Fatal(err)
	}
	client.Clear()
	defer client.Clear()

	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}

	valueBefore := "value"
	client.Set("key", &valueBefore)
	var value string
	if err := client.Get("key", &value); err != nil || value != "value" {
		t.Error("Expected value to be 'value'")
	}

	if client.Size() != 1 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}
}
