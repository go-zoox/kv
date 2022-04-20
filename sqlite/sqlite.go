package sqlite

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLite is a Key-Value Store in SQLite
type SQLite struct {
	sync.RWMutex
	Core   *sql.DB
	Config *SQLiteConfig
}

// SQLiteConfig is the configuration for Redis
type SQLiteConfig struct {
	// Path to the SQLite database file.
	Path string

	// Prefix is the prefix to use for all keys
	Prefix string
}

// New returns a new MemoryKV.
func New(cfg *SQLiteConfig) (*SQLite, error) {
	if cfg.Path == "" {
		return nil, errors.New("sqlite: path is required")
	}

	if cfg.Prefix == "" {
		return nil, errors.New("prefix is required")
	}

	core, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, err
	}

	// Create the table if it doesn't exist
	_, err = core.Exec("CREATE TABLE IF NOT EXISTS kv (key TEXT PRIMARY KEY, value BLOB, expires_at INTEGER)")
	if err != nil {
		return nil, err
	}

	return &SQLite{
		Core:   core,
		Config: cfg,
	}, nil
}

func (m *SQLite) getKey(key string) string {
	return m.Config.Prefix + key
}

func now() int64 {
	return time.Now().Unix() * 1000
}

// Set sets the value for the given key.
// If maxAge is greater than 0, then the value will be expired after maxAge miliseconds.
func (m *SQLite) Set(key string, value string, maxAge ...int64) error {
	m.Lock()
	defer m.Unlock()

	var maxAgeX int64
	var expiresAt int64
	if len(maxAge) > 0 {
		maxAgeX = int64(maxAge[0])
		expiresAt = now() + maxAgeX
	}

	keyX := m.getKey(key)
	stmt, err := m.Core.Prepare("INSERT OR REPLACE INTO kv (key, value, expires_at) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(keyX, value, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

// Get returns the value for the given key.
func (m *SQLite) Get(key string) string {
	m.RLock()

	keyX := m.getKey(key)
	stmt, err := m.Core.Prepare("SELECT value, expires_at FROM kv WHERE key = ?")
	if err != nil {
		panic(err)
	}

	res := stmt.QueryRow(keyX)
	if res.Err() != nil {
		panic(res.Err())
	}

	var value string
	var expiresAt int64
	if err := res.Scan(&value, &expiresAt); err != nil {
		// panic(err)
		m.RUnlock()
		return ""
	}

	m.RUnlock()
	if expiresAt > 0 && expiresAt < now() {
		m.Delete(key)
		return ""
	}

	return value
}

// Delete deletes the value for the given key.
func (m *SQLite) Delete(key string) error {
	m.Lock()
	defer m.Unlock()

	stmt, err := m.Core.Prepare("DELETE FROM kv WHERE key = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(m.getKey(key))
	if err != nil {
		return err
	}

	return nil
}

// Has returns true if the given key exists in the kv.
func (m *SQLite) Has(key string) bool {
	m.RLock()
	defer m.RUnlock()

	stmt, err := m.Core.Prepare("SELECT 1 FROM kv WHERE key = ?")
	if err != nil {
		panic(err)
	}

	res := stmt.QueryRow(m.getKey(key))
	if res.Err() != nil {
		panic(res.Err())
	}

	var value int
	if err := res.Scan(&value); err != nil {
		return false
	}

	return value > 0
}

// Keys returns the keys of the kv.
func (m *SQLite) Keys() []string {
	m.RLock()
	defer m.RUnlock()

	rows, err := m.Core.Query("SELECT key FROM kv where key like ?", m.Config.Prefix+"%")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	keys := make([]string, 0)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			panic(err)
		}

		keys = append(keys, key[len(m.Config.Prefix):])
	}

	return keys
}

// Size returns the number of elements in the kv.
func (m *SQLite) Size() int {
	m.RLock()
	defer m.RUnlock()

	res, err := m.Core.Query("SELECT count(*) FROM kv where key like ?", m.Config.Prefix+"%")
	if err != nil {
		panic(err)
	}
	defer res.Close()

	var count int
	res.Next()
	if err := res.Scan(&count); err != nil {
		panic(err)
	}

	return count
}

// Clear removes all elements from the kv.
func (m *SQLite) Clear() error {
	m.Lock()
	defer m.Unlock()

	_, err := m.Core.Exec("DELETE FROM kv where key like ?", m.Config.Prefix+"%")
	return err
}

// ForEach calls the given function for each key-value pair in the kv.
func (m *SQLite) ForEach(f func(string, interface{})) {
	m.RLock()
	defer m.RUnlock()

	keys := m.Keys()
	for _, key := range keys {
		f(key, m.Get(key))
	}
}
