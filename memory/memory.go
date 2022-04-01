package memory

import (
	"sync"
	"time"
)

// Memory is a Key-Value Store in Memory，like JavaScript Map for Go
type Memory struct {
	sync.RWMutex
	data map[string]interface{}
}

// Value is a value of Memory
type Value struct {
	Value     interface{}
	ExpiresAt int64
}

// New returns a new MemoryKV.
func New() *Memory {
	return &Memory{
		data: make(map[string]interface{}),
	}
}

func now() int64 {
	return time.Now().Unix() * 1000
}

// Set sets the value for the given key.
// If maxAge is greater than 0, then the value will be expired after maxAge miliseconds.
func (m *Memory) Set(key string, value interface{}, maxAge ...int64) error {
	m.Lock()
	defer m.Unlock()

	expiresAt := int64(0)
	if len(maxAge) > 0 {
		expiresAt = now() + maxAge[0]
	}

	m.data[key] = Value{value, expiresAt}
	return nil
}

// Get returns the value for the given key.
func (m *Memory) Get(key string) interface{} {
	m.RLock()
	defer m.RUnlock()

	if !m.Has(key) {
		return nil
	}

	v := m.data[key].(Value)
	if v.ExpiresAt > 0 && v.ExpiresAt < now() {
		delete(m.data, key)
		return nil
	}

	return v.Value
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
	defer m.RUnlock()

	_, ok := m.data[key]
	if !ok {
		return false
	}

	v := m.data[key].(Value)
	if v.ExpiresAt > 0 && v.ExpiresAt < now() {
		delete(m.data, key)
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

	for k := range m.data {
		delete(m.data, k)
	}
	return nil
}

// ForEach calls the given function for each key-value pair in the kv.
func (m *Memory) ForEach(f func(string, interface{})) {
	m.RLock()
	defer m.RUnlock()

	for k, v := range m.data {
		f(k, v.(Value).Value)
	}
}
