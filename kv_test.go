package kv

import "testing"

func TestKVMemory(t *testing.T) {
	m := NewMemory()
	if m.Size() != 0 {
		t.Errorf("Expected size 0, got %d", m.Size())
	}

	m.Set("key", "value")
	if m.Get("key") != "value" {
		t.Error("Expected value to be 'value'")
	}

	if m.Size() != 1 {
		t.Errorf("Expected size 1, got %d", m.Size())
	}
}
