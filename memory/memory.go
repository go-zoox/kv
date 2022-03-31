package memory

// MemoryKV is a Memory Key-Value Store，like JavaScript Map for Go
type MemoryKV map[string]interface{}

// MemoryKVEntry is a key-value pair.
type MemoryKVEntry struct {
	Key   string
	Value interface{}
}

// Set sets the value for the given key.
func (m MemoryKV) Set(key string, value interface{}) error {
	m[key] = value
	return nil
}

// Get returns the value for the given key.
func (m MemoryKV) Get(key string) interface{} {
	return m[key]
}

// Delete deletes the value for the given key.
func (m MemoryKV) Delete(key string) error {
	delete(m, key)
	return nil
}

// Has returns true if the given key exists in the kv.
func (m MemoryKV) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// Keys returns the keys of the kv.
func (m MemoryKV) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Values returns the values of the kv.
func (m MemoryKV) Values() []interface{} {
	values := make([]interface{}, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

// Entries is an iterator for MapEntry.
func (m MemoryKV) Entries() []MemoryKVEntry {
	entries := make([]MemoryKVEntry, len(m))
	i := 0
	for k, v := range m {
		entries[i] = MemoryKVEntry{k, v}
		i++
	}
	return entries
}

// Size returns the number of elements in the kv.
func (m MemoryKV) Size() int {
	return len(m)
}

// Clear removes all elements from the kv.
func (m MemoryKV) Clear() error {
	for k := range m {
		delete(m, k)
	}
	return nil
}

// ForEach calls the given function for each key-value pair in the kv.
func (m MemoryKV) ForEach(f func(string, interface{})) {
	for k, v := range m {
		f(k, v)
	}
}
