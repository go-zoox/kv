package test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-zoox/kv/typing"
)

// RunTestCases runs all test cases.
func RunTestCases(t *testing.T, client typing.KV, casesDisabled ...[]string) {
	client.Clear()
	defer client.Clear()

	casesDisabledX := map[string]bool{
		"main":    true,
		"keys":    true,
		"forEach": true,
		"maxAge":  true,
	}
	if len(casesDisabled) > 0 {
		for _, c := range casesDisabled[0] {
			casesDisabledX[c] = false
		}
	}

	if casesDisabledX["main"] {
		RunMainTestCase(t, client)
	}

	if casesDisabledX["keys"] {
		RunKeysTestCase(t, client)
	}

	if casesDisabledX["forEach"] {
		RunForEachTestCase(t, client)
	}

	if casesDisabledX["maxAge"] {
		RunMaxAgeTestCase(t, client)
	}
}

// RunMainTestCase tests the main functionality.
func RunMainTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing main test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}
	if !client.Has("key") {
		t.Error("Expected key to be set")
	}
	if client.Size() != 1 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}

	if err := client.Delete("key"); err != nil {
		t.Fatal(err)
	}
	if client.Has("key") {
		t.Error("Expected key to be deleted")
	}
	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 1 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 2 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}
}

// RunKeysTestCase tests the keys functionality.
func RunKeysTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing keys test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key3", "value3"); err != nil {
		t.Fatal(err)
	}

	keys := client.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected len 3, got %d", len(keys))
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})
	if strings.Join(keys, ",") != "key1,key2,key3" {
		t.Errorf("Expected keys to be key1,key2,key3, got %v", strings.Join(keys, ","))
	}
}

// RunForEachTestCase tests the ForEach functionality.
func RunForEachTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing forEach test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key3", "value3"); err != nil {
		t.Fatal(err)
	}

	// client.ForEach(func(key string, value interface{}) {
	// 	if key == "key1" && value != "value1" {
	// 		t.Errorf("Expected value to be 'value1', but got %s", value)
	// 	}
	// 	if key == "key2" && value != "value2" {
	// 		t.Errorf("Expected value to be 'value2', but got %s", value)
	// 	}
	// 	if key == "key3" && value != "value3" {
	// 		t.Errorf("Expected value to be 'value3', but got %s", value)
	// 	}
	// })
}

// RunMaxAgeTestCase tests the maxAge functionality.
func RunMaxAgeTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing max age test case")

	client.Clear()
	defer client.Clear()

	key1 := "key1"
	value1 := "value1"
	if err := client.Set(key1, &value1, 500); err != nil {
		t.Fatal(err)
	}

	var value string
	if err := client.Get("key1", &value); err != nil || value != value1 {
		t.Error("Expected value to be 'value1'")
	}

	// fmt.Println("xxx1:", client.Get("key1"))
	time.Sleep(2 * time.Second)
	// done := make(chan bool)
	// go func() {
	// 	fmt.Println("xxx2:", client.Get("key1"))
	// 	time.Sleep(2 * time.Second)
	// 	fmt.Println("xxx3:", client.Get("key1"))
	// 	done <- true
	// }()
	// <-done
	// fmt.Println("xxx4:", client.Get("key1"))

	client.Get("key1", &value)
	if client.Has("key1") {
		t.Errorf("Expected value to be '', but got %s", value)
	}
}
