package memory

import (
	"testing"
	"time"

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

func TestPlainStringGetSet(t *testing.T) {
	client := createClient()
	defer client.Clear()

	if err := client.Set("k", "string value"); err != nil {
		t.Fatal(err)
	}
	var v string
	if err := client.Get("k", &v); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if v != "string value" {
		t.Fatalf("expected 'string value', got %q", v)
	}
}

func TestPlainIntGetSet(t *testing.T) {
	client := createClient()
	defer client.Clear()

	if err := client.Set("k", 42); err != nil {
		t.Fatal(err)
	}
	var v int
	if err := client.Get("k", &v); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestPlainBoolGetSet(t *testing.T) {
	client := createClient()
	defer client.Clear()

	if err := client.Set("k", true); err != nil {
		t.Fatal(err)
	}
	var v bool
	if err := client.Get("k", &v); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !v {
		t.Fatal("expected true")
	}
}

func TestPlainStringSliceGetSet(t *testing.T) {
	client := createClient()
	defer client.Clear()

	orig := []string{"x", "y", "z"}
	if err := client.Set("k", orig); err != nil {
		t.Fatal(err)
	}
	var v []string
	if err := client.Get("k", &v); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(v) != 3 || v[0] != "x" || v[2] != "z" {
		t.Fatalf("expected [x y z], got %v", v)
	}
}

func TestGetNonexistentKey(t *testing.T) {
	client := createClient()
	defer client.Clear()

	var v string
	err := client.Get("no-such", &v)
	if err == nil {
		t.Fatal("expected error for nonexistent key")
	}
}

func TestHasAfterExpiry(t *testing.T) {
	client := createClient()
	defer client.Clear()

	if err := client.Set("k", "v", 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if !client.Has("k") {
		t.Fatal("expected key to exist")
	}

	time.Sleep(150 * time.Millisecond)
	if client.Has("k") {
		t.Fatal("expected key to expire")
	}
}

func TestSize(t *testing.T) {
	client := createClient()
	defer client.Clear()

	if client.Size() != 0 {
		t.Fatalf("expected size 0, got %d", client.Size())
	}
	client.Set("a", 1)
	client.Set("b", 2)
	if client.Size() != 2 {
		t.Fatalf("expected size 2, got %d", client.Size())
	}
	client.Delete("a")
	if client.Size() != 1 {
		t.Fatalf("expected size 1, got %d", client.Size())
	}
}

func TestKeys(t *testing.T) {
	client := createClient()
	defer client.Clear()

	client.Set("b", 2)
	client.Set("a", 1)
	client.Set("c", 3)

	keys := client.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
}

func TestClear(t *testing.T) {
	client := createClient()
	defer client.Clear()

	client.Set("a", 1)
	client.Set("b", 2)
	client.Clear()
	if client.Size() != 0 {
		t.Fatalf("expected size 0 after clear, got %d", client.Size())
	}
	if client.Has("a") || client.Has("b") {
		t.Fatal("expected no keys after clear")
	}
}

func TestOverwriteValue(t *testing.T) {
	client := createClient()
	defer client.Clear()

	client.Set("k", "first")
	client.Set("k", "second")

	var v string
	client.Get("k", &v)
	if v != "second" {
		t.Fatalf("expected 'second', got %q", v)
	}
	if client.Size() != 1 {
		t.Fatalf("expected size 1, got %d", client.Size())
	}
}

func TestSetWithTTL(t *testing.T) {
	client := createClient()
	defer client.Clear()

	if err := client.Set("k", "v", 5*time.Second); err != nil {
		t.Fatal(err)
	}
	if !client.Has("k") {
		t.Fatal("expected key to exist")
	}
	var v string
	if err := client.Get("k", &v); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if v != "v" {
		t.Fatalf("expected 'v', got %q", v)
	}
}
