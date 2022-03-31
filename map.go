package jmap

// Map is JavaScript Map Like for Go
type Map map[string]interface{}

// MapEntry is a key-value pair.
type MapEntry struct {
	Key   string
	Value interface{}
}

// Set sets the value for the given key.
func (m Map) Set(key string, value interface{}) {
	m[key] = value
}

// Get returns the value for the given key.
func (m Map) Get(key string) interface{} {
	return m[key]
}

// Delete deletes the value for the given key.
func (m Map) Delete(key string) {
	delete(m, key)
}

// Has returns true if the given key exists in the map.
func (m Map) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// Keys returns the keys of the map.
func (m Map) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Values returns the values of the map.
func (m Map) Values() []interface{} {
	values := make([]interface{}, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

// Entries is an iterator for MapEntry.
func (m Map) Entries() []MapEntry {
	entries := make([]MapEntry, len(m))
	i := 0
	for k, v := range m {
		entries[i] = MapEntry{k, v}
		i++
	}
	return entries
}

// Size returns the number of elements in the map.
func (m Map) Size() int {
	return len(m)
}

// Clear removes all elements from the map.
func (m Map) Clear() {
	for k := range m {
		delete(m, k)
	}
}

// ForEach calls the given function for each key-value pair in the map.
func (m Map) ForEach(f func(string, interface{})) {
	for k, v := range m {
		f(k, v)
	}
}
