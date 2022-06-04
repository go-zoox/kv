package fs

import (
	"os"
	"reflect"
	"sync"
	"time"

	zfs "github.com/go-zoox/fs"
	zjson "github.com/go-zoox/fs/type/json"
)

// FileSystem is a Key-Value Store in FileSystem，like JavaScript Map for Go
type FileSystem struct {
	sync.RWMutex
	dir string
}

// Value is a value of Memory
type Value struct {
	Value     any
	ExpiresAt int64
}

// FileSystemOptions represents the options for the kv.
type FileSystemOptions struct {
	Dir string
}

// New returns a new MemoryKV.
func New(cfg ...*FileSystemOptions) (*FileSystem, error) {
	homeDir, _ := os.UserHomeDir()
	dir := zfs.JoinPath(homeDir, ".cache/go-zoox/kv/fs")
	if len(cfg) > 0 && cfg[0] != nil {
		if cfg[0].Dir != "" {
			dir = cfg[0].Dir
		}
	}

	return &FileSystem{
		dir: dir,
	}, nil
}

func now() int64 {
	return time.Now().UnixMilli()
}

// Set sets the value for the given key.
// If maxAge is greater than 0, then the value will be expired after maxAge miliseconds.
func (m *FileSystem) Set(key string, value any, maxAge ...time.Duration) error {
	m.Lock()
	// defer m.Unlock()

	m.ensureDir()

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
			v := m.read(key)
			expiresAt = v.ExpiresAt
		}

		m.Lock()
	}

	if err := m.write(key, &Value{value, expiresAt}); err != nil {
		m.Unlock()
		return err
	}

	m.Unlock()
	return nil
}

// Get returns the value for the given key.
func (m *FileSystem) Get(key string, value interface{}) error {
	m.RLock()

	if !m.Has(key) {
		m.RUnlock()
		return nil
	}

	v := m.read(key) // m.data[key].(Value)
	m.RUnlock()

	if v.ExpiresAt > 0 && v.ExpiresAt < now() {
		m.Delete(key)
		return nil
	}

	// reflect valueof
	reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v.Value))
	return nil
}

// Delete deletes the value for the given key.
func (m *FileSystem) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	return m.remove(key)
}

// Has returns true if the given key exists in the kv.
func (m *FileSystem) Has(key string) bool {
	m.RLock()
	defer m.RUnlock()

	v := m.read(key)
	if v == nil {
		return false
	}

	if v.ExpiresAt > 0 && v.ExpiresAt < now() {
		m.remove(key)
		return false
	}

	return true
}

// Keys returns the keys of the kv.
func (m *FileSystem) Keys() []string {
	m.RLock()
	defer m.RUnlock()

	files, err := zfs.ListDir(m.dir)
	if err != nil {
		return []string{}
	}

	keys := make([]string, len(files))
	i := 0
	for _, f := range files {
		keys[i] = f.Name()
		i++
	}
	return keys
}

// Size returns the number of elements in the kv.
func (m *FileSystem) Size() int {
	m.RLock()
	defer m.RUnlock()

	files, err := zfs.ListDir(m.dir)
	if err != nil {
		return 0
	}

	return len(files)
}

// Clear removes all elements from the kv.
func (m *FileSystem) Clear() error {
	m.Lock()
	defer m.Unlock()

	return zfs.RemoveDir(m.dir)
}

// ForEach calls the given function for each key-value pair in the kv.
func (m *FileSystem) ForEach(f func(string, interface{})) {
	m.RLock()
	files, err := zfs.ListDir(m.dir)
	m.RUnlock()

	if err != nil {
		return
	}

	for _, file := range files {
		k := file.Name()
		v := m.read(k)
		if v != nil {
			f(k, v.Value)
		}
	}
}

func (m *FileSystem) ensureDir() {
	if zfs.IsExist(m.dir) {
		return
	}

	zfs.Mkdirp(m.dir)
}

func (m *FileSystem) filepathOfKey(key string) string {
	return zfs.JoinPath(m.dir, key)
}

func (m *FileSystem) write(key string, v *Value) error {
	filepath := m.filepathOfKey(key)
	if !zfs.IsExist(filepath) {
		zfs.CreateFile(filepath)
	}

	// fmt.Println("write file:", filepath, v)
	return zjson.Write(filepath, v)
}

func (m *FileSystem) read(key string) *Value {
	filepath := m.filepathOfKey(key)
	if !zfs.IsExist(filepath) {
		return nil
	}

	var v Value
	if err := zjson.Read(filepath, &v); err != nil {
		return nil
	}

	// fmt.Println("key:", key, "value:", v)
	return &v
}

func (m *FileSystem) remove(key string) error {
	filepath := m.filepathOfKey(key)
	if !zfs.IsExist(filepath) {
		return nil
	}

	return zfs.Remove(filepath)
}
