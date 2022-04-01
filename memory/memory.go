package memory

// Memory is a Key-Value Store in Memory，like JavaScript Map for Go
type Memory map[string]interface{}

// New returns a new MemoryKV.
func New() *Memory {
	return &Memory{}
}

// Set sets the value for the given key.
func (m Memory) Set(key string, value interface{}) error {
	m[key] = value
	return nil
}

// Get returns the value for the given key.
func (m Memory) Get(key string) interface{} {
	return m[key]
}

// Delete deletes the value for the given key.
func (m Memory) Delete(key string) error {
	delete(m, key)
	return nil
}

// Has returns true if the given key exists in the kv.
func (m Memory) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// Keys returns the keys of the kv.
func (m Memory) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Values returns the values of the kv.
func (m Memory) Values() []interface{} {
	values := make([]interface{}, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

// Size returns the number of elements in the kv.
func (m Memory) Size() int {
	return len(m)
}

// Clear removes all elements from the kv.
func (m Memory) Clear() error {
	for k := range m {
		delete(m, k)
	}
	return nil
}

// ForEach calls the given function for each key-value pair in the kv.
func (m Memory) ForEach(f func(string, interface{})) {
	for k, v := range m {
		f(k, v)
	}
}
