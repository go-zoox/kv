package typing

// KV is a Key-Value Store
type KV interface {
	// Set sets the value for the given key.
	Set(key string, value interface{}) error
	// Get returns the value for the given key.
	Get(key string) interface{}
	// Delete deletes the value for the given key.
	Delete(key string) error
	// Has returns true if the given key exists in the kv.
	Has(key string) bool
	// Keys returns the keys of the kv.
	Keys() []string
	// Values returns the values of the kv.
	Values() []interface{}
	// Size returns the number of entries in the kv.
	Size() int
	// Clear clears the kv.
	Clear() error
	// ForEach iterates over the map and calls the given function for each entry.
	ForEach(func(key string, value interface{}))
}
