package memory

import (
	"reflect"
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
	Value     any
	ExpiresAt int64
}

// New returns a new MemoryKV.
func New() *Memory {
	return &Memory{
		data: make(map[string]interface{}),
	}
}

func now() int64 {
	return time.Now().UnixMilli()
}

// Set sets the value for the given key.
// If maxAge is greater than 0, then the value will be expired after maxAge miliseconds.
func (m *Memory) Set(key string, value any, maxAge ...time.Duration) error {
	m.Lock()
	// defer m.Unlock()

	expiresAt := int64(0)
	if len(maxAge) > 0 {
		expiresAt = now() + int64(maxAge[0]/time.Millisecond)
	} else {
		m.Unlock()

		if m.Has(key) {
			// var v Value
			// if err := m.Get(key, &v); err != nil {
			// 	return err
			// }

			// use origin expiresAt
			v := m.data[key].(Value)
			expiresAt = v.ExpiresAt
		}

		m.Lock()
	}

	m.data[key] = Value{value, expiresAt}
	m.Unlock()
	return nil
}

// Get returns the value for the given key.
func (m *Memory) Get(key string, value interface{}) error {
	m.RLock()

	if !m.Has(key) {
		m.RUnlock()
		return nil
	}

	v := m.data[key].(Value)
	m.RUnlock()

	if v.ExpiresAt > 0 && v.ExpiresAt < now() {
		m.Delete(key)
		return nil
	}

	reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v.Value))
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

	_, ok := m.data[key]
	if !ok {
		m.RUnlock()
		return false
	}

	v := m.data[key].(Value)
	m.RUnlock()

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
