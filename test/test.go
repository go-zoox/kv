package test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-zoox/kv/typing"
)

func RunTestCases(t *testing.T, client typing.KV) {
	client.Clear()
	defer client.Clear()

	RunMainTestCase(t, client)
	RunKeysTestCase(t, client)
	RunForEachTestCase(t, client)
	RunMaxAgeTestCase(t, client)
}

func RunMainTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing main test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}
	if !client.Has("key") {
		t.Error("Expected key to be set")
	}
	if client.Size() != 1 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}

	if err := client.Delete("key"); err != nil {
		t.Fatal(err)
	}
	if client.Has("key") {
		t.Error("Expected key to be deleted")
	}
	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 1 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 2 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}
}

func RunKeysTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing keys test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key3", "value3"); err != nil {
		t.Fatal(err)
	}

	keys := client.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected len 3, got %d", len(keys))
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})
	if strings.Join(keys, ",") != "key1,key2,key3" {
		t.Errorf("Expected keys to be key1,key2,key3, got %v", strings.Join(keys, ","))
	}
}

func RunForEachTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing forEach test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key3", "value3"); err != nil {
		t.Fatal(err)
	}

	client.ForEach(func(key string, value interface{}) {
		if key == "key1" && value != "value1" {
			t.Error("Expected value to be 'value1'")
		}
		if key == "key2" && value != "value2" {
			t.Error("Expected value to be 'value2'")
		}
		if key == "key3" && value != "value3" {
			t.Error("Expected value to be 'value3'")
		}
	})
}

func RunMaxAgeTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing max age test case")

	client.Clear()
	defer client.Clear()

	client.Set("key1", "value1", 500)
	if client.Get("key1") != "value1" {
		t.Error("Expected value to be 'value1'")
	}

	time.Sleep(2 * time.Second)
	if client.Has("key1") {
		t.Errorf("Expected value to be '', but got %s", client.Get("key1"))
	}
}
