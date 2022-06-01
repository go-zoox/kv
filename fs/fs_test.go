package fs

import (
	"testing"

	"github.com/go-zoox/kv/test"
)

func createClient() *FileSystem {
	client := New()
	return client
}

func TestKV(t *testing.T) {
	test.RunTestCases(t, createClient())
	// client := createClient()
	// client.Set("key", "777")
	// var v string
	// client.Get("key", &v)
	// fmt.Println("v:", v)
}
