package memory

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Memory is a Key-Value Store in Memory, like JavaScript Map for Go.
type Memory struct {
	sync.RWMutex
	data map[string]Value
}

// Value is a value of Memory.
type Value struct {
	Value     interface{}
	ExpiresAt int64
}

// New returns a new MemoryKV.
func New() *Memory {
	return &Memory{
		data: make(map[string]Value),
	}
}

func now() int64 {
	return time.Now().UnixMilli()
}

// Set sets the value for the given key.
// If maxAge is greater than 0, the value will be expired after maxAge milliseconds.
func (m *Memory) Set(key string, value interface{}, maxAge ...time.Duration) error {
	m.Lock()
	defer m.Unlock()

	if value == nil {
		return fmt.Errorf("value is nil")
	}

	expiresAt := int64(0)
	if len(maxAge) > 0 {
		expiresAt = now() + int64(maxAge[0]/time.Millisecond)
	} else if val, ok := m.data[key]; ok {
		expiresAt = val.ExpiresAt
	}

	m.data[key] = Value{value, expiresAt}
	return nil
}

// Get returns the value for the given key.
func (m *Memory) Get(key string, value interface{}) error {
	m.RLock()
	val, ok := m.data[key]
	m.RUnlock()

	if !ok {
		return fmt.Errorf("key %s not found", key)
	}

	if val.ExpiresAt > 0 && val.ExpiresAt < now() {
		m.Delete(key)
		return fmt.Errorf("key %s expired", key)
	}

	// reference: https://riptutorial.com/go/example/6073/reflect-value-elem--
	v := reflect.ValueOf(val.Value)
	if v.IsZero() {
		return nil
	}

	reflect.ValueOf(value).Elem().Set(v.Elem())
	return nil
}

// Delete deletes the value for the given key.
func (m *Memory) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.data, key)
	return nil
}

// Has returns true if the given key exists in the kv.
func (m *Memory) Has(key string) bool {
	m.RLock()
	val, ok := m.data[key]
	m.RUnlock()

	if !ok {
		return false
	}

	if val.ExpiresAt > 0 && val.ExpiresAt < now() {
		m.Delete(key)
		return false
	}

	return true
}

// Keys returns the keys of the kv.
func (m *Memory) Keys() []string {
	m.RLock()
	defer m.RUnlock()

	keys := make([]string, len(m.data))
	i := 0
	for k := range m.data {
		keys[i] = k
		i++
	}

	return keys
}

// Size returns the number of elements in the kv.
func (m *Memory) Size() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.data)
}

// Clear removes all elements from the kv.
func (m *Memory) Clear() error {
	m.Lock()
	defer m.Unlock()

	m.data = make(map[string]Value)
	return nil
}

// ForEach calls the given function for each key-value pair in the kv.
func (m *Memory) ForEach(f func(string, interface{})) {
	m.RLock()
	defer m.RUnlock()

	for k, v := range m.data {
		f(k, v.Value)
	}
}
