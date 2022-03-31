package memory

import (
	"sort"
	"strings"
	"testing"
)

func TestMapGetSet(t *testing.T) {
	m := Memory{}
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

func TestMapDeleteClear(t *testing.T) {
	m := Memory{}
	m.Set("key", "value")
	m.Delete("key")
	if m.Has("key") {
		t.Error("Expected key to be deleted")
	}
	m.Set("key", "value")
	m.Clear()
	if m.Size() != 0 {
		t.Errorf("Expected size 0, got %d", m.Size())
	}
	if m.Has("key") {
		t.Error("Expected key to be deleted")
	}
}

func TestMapKeysValues(t *testing.T) {
	m := Memory{}
	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")
	keys := m.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected len 3, got %d", len(keys))
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})
	if strings.Join(keys, ",") != "key1,key2,key3" {
		t.Errorf("Expected keys to be key1,key2,key3, got %v", strings.Join(keys, ","))
	}
	values := m.Values()
	if len(values) != 3 {
		t.Errorf("Expected len 3, got %d", len(values))
	}
	sort.Slice(values, func(i, j int) bool {
		return strings.Compare(values[i].(string), values[j].(string)) < 0
	})
	valuesX := []string{}
	for _, v := range values {
		valuesX = append(valuesX, v.(string))
	}
	if strings.Join(valuesX, ",") != "value1,value2,value3" {
		t.Errorf("Expected values to be value1,value2,value3, got %v", strings.Join(valuesX, ","))
	}
}

func TestMapForEach(t *testing.T) {
	m := Memory{}
	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")
	m.ForEach(func(key string, value interface{}) {
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
